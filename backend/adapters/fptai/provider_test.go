package fptai

import (
	"context"
	"errors"
	"testing"
)

func TestUnavailableProvidersReturnSanitizedError(t *testing.T) {
	if _, err := NewUnavailableASR().Transcribe(context.Background(), AudioInput{}); !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("expected unavailable ASR error, got %v", err)
	}
	if _, err := NewUnavailableTranslation().Translate(context.Background(), TranslationInput{}); !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("expected unavailable translation error, got %v", err)
	}
}

func TestUnimplementedRealCapabilitiesDoNotReturnSimulatedSuccess(t *testing.T) {
	client := NewClient("https://provider.invalid", "test-key", 0)

	if result, err := NewRealASR(client, "asr").Transcribe(context.Background(), AudioInput{}); !errors.Is(err, ErrCapabilityUnavailable) || result.Text != "" {
		t.Fatalf("real ASR must explicitly degrade, result=%+v err=%v", result, err)
	}
	if result, err := NewRealVision(client, "vision").Observe(context.Background(), VisionInput{}); !errors.Is(err, ErrCapabilityUnavailable) || result.Description != "" {
		t.Fatalf("real vision must explicitly degrade, result=%+v err=%v", result, err)
	}
	if result, err := NewRealExtraction(client, "extraction").Extract(context.Background(), ExtractionInput{}); !errors.Is(err, ErrCapabilityUnavailable) || len(result.Facts) != 0 {
		t.Fatalf("real extraction must explicitly degrade, result=%+v err=%v", result, err)
	}
}
