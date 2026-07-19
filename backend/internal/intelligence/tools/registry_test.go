package tools

import (
	"context"
	"testing"
)

func TestRegistryRejectsArbitraryTool(t *testing.T) {
	result := NewRegistry().Execute(context.Background(), "run_sql", []byte(`{}`), "trace")
	if result.Status != "failed" || result.ErrorCategory != "tool_not_allowed" {
		t.Fatalf("arbitrary tool was not rejected: %+v", result)
	}
}
