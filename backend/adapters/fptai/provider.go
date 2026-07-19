package fptai

import (
	"context"
	"errors"
)

var (
	// ErrProviderUnavailable is intentionally generic so callers can safely
	// surface a degraded state without leaking provider configuration details.
	ErrProviderUnavailable = errors.New("AI provider unavailable")
	// ErrCapabilityUnavailable indicates that a configured adapter does not yet
	// implement the requested capability. It must never be replaced by mock data.
	ErrCapabilityUnavailable = errors.New("AI capability unavailable")
)

type AudioInput struct {
	Data []byte
}

type Transcript struct {
	Text string
}

type TranslationInput struct {
	Text   string
	Target string // Target language locale
}

type Translation struct {
	Text string
}

type VisionInput struct {
	ImageBytes []byte
}

type VisionObservation struct {
	Description string
}

type ExtractionInput struct {
	Text string
}

type StructuredFacts struct {
	Facts []string
}

type ASRProvider interface {
	Transcribe(ctx context.Context, input AudioInput) (Transcript, error)
}

type TranslationProvider interface {
	Translate(ctx context.Context, input TranslationInput) (Translation, error)
}

type VisionProvider interface {
	Observe(ctx context.Context, input VisionInput) (VisionObservation, error)
}

type ExtractionProvider interface {
	Extract(ctx context.Context, input ExtractionInput) (StructuredFacts, error)
}

type UnavailableASR struct{}

func NewUnavailableASR(_ ...string) *UnavailableASR {
	return &UnavailableASR{}
}

func (p *UnavailableASR) Transcribe(context.Context, AudioInput) (Transcript, error) {
	return Transcript{}, ErrProviderUnavailable
}

type UnavailableTranslation struct{}

func NewUnavailableTranslation(_ ...string) *UnavailableTranslation {
	return &UnavailableTranslation{}
}

func (p *UnavailableTranslation) Translate(context.Context, TranslationInput) (Translation, error) {
	return Translation{}, ErrProviderUnavailable
}
