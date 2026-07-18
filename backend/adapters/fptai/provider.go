package fptai

import (
	"context"
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
