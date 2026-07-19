package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tourtect/backend/internal/intelligence/model"
)

var ErrToolNotAllowed = errors.New("tool is not allowlisted")

type Spec struct {
	Name                 string
	Description          string
	InputSchema          string
	OutputSchema         string
	RequiredConsent      string
	RequiresConfirmation bool
	Timeout              time.Duration
	ErrorBehavior        string
	AuditBehavior        string
}

type Tool interface {
	Spec() Spec
	Execute(context.Context, json.RawMessage, string) (json.RawMessage, string, error)
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry(allowed ...Tool) *Registry {
	r := &Registry{tools: make(map[string]Tool, len(allowed))}
	for _, tool := range allowed {
		r.tools[tool.Spec().Name] = tool
	}
	return r
}

func (r *Registry) Specs() []Spec {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Spec, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool.Spec())
	}
	return result
}

func (r *Registry) Execute(ctx context.Context, name string, input json.RawMessage, traceID string) model.ToolResult {
	r.mu.RLock()
	tool, ok := r.tools[name]
	r.mu.RUnlock()
	started := time.Now()
	result := model.ToolResult{ID: uuid.NewString(), ToolName: name, Status: "failed", Output: json.RawMessage(`{}`)}
	if !ok {
		result.ErrorCategory = "tool_not_allowed"
		return result
	}
	timeout := tool.Spec().Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	toolCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	output, status, err := tool.Execute(toolCtx, input, traceID)
	result.DurationMS = time.Since(started).Milliseconds()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			result.ErrorCategory = "timeout"
		} else {
			result.ErrorCategory = "tool_error"
		}
		return result
	}
	if !json.Valid(output) {
		result.ErrorCategory = "invalid_tool_output"
		return result
	}
	result.Output, result.Status = output, status
	return result
}

func DecodeInput[T any](input json.RawMessage) (T, error) {
	var value T
	if err := json.Unmarshal(input, &value); err != nil {
		return value, fmt.Errorf("invalid tool input: %w", err)
	}
	return value, nil
}
