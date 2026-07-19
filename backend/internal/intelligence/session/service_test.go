package session

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/tourtect/backend/internal/intelligence/model"
)

func TestSessionOwnershipAndBoundedFacts(t *testing.T) {
	ctx := context.Background()
	service := NewService(NewMemoryStore(), time.Minute)
	value, err := service.Create(ctx, "user-a", "en", "vi-VN", "", "hanoi", "text", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetOwned(ctx, value.ID, "user-b"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ownership boundary, got %v", err)
	}
	for i := 0; i < 40; i++ {
		value.Context.UserConfirmedFacts = append(value.Context.UserConfirmedFacts, strings.Repeat("x", 700))
	}
	if err := service.Save(ctx, value); err != nil {
		t.Fatal(err)
	}
	stored, err := service.GetOwned(ctx, value.ID, "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(stored.Context.UserConfirmedFacts) != 32 {
		t.Fatalf("facts not bounded: %d", len(stored.Context.UserConfirmedFacts))
	}
	if len([]rune(stored.Context.UserConfirmedFacts[0])) != 512 {
		t.Fatal("fact was not redacted to maximum length")
	}
}

func TestMessageIdempotency(t *testing.T) {
	ctx := context.Background()
	service := NewService(NewMemoryStore(), time.Minute)
	value, _ := service.Create(ctx, "u", "en", "", "", "", "text", false)
	response := model.Response{ID: "response", Intent: "unknown", Evidence: []model.Evidence{}, ToolResults: []model.ToolResult{}, SuggestedActions: []model.SuggestedAction{}, FactorsConsidered: []string{}, MissingInformation: []string{}}
	if err := service.AppendResponse(ctx, value, "message", response); err != nil {
		t.Fatal(err)
	}
	if !service.HasMessage(value, "message") {
		t.Fatal("processed message id not retained")
	}
}
