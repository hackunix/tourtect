package fptai

import (
	"context"
	"errors"
)

type FakeASR struct {
	ShouldFail bool
	MockText   string
}

func (f *FakeASR) Transcribe(ctx context.Context, input AudioInput) (Transcript, error) {
	if f.ShouldFail {
		return Transcript{}, errors.New("fake ASR error")
	}
	text := f.MockText
	if text == "" {
		text = "Chào anh, taxi từ Nội Bài về Hoàn Kiếm giá bao nhiêu?"
	}
	return Transcript{Text: text}, nil
}

type FakeTranslation struct {
	ShouldFail bool
	MockText   string
}

func (f *FakeTranslation) Translate(ctx context.Context, input TranslationInput) (Translation, error) {
	if f.ShouldFail {
		return Translation{}, errors.New("fake translation error")
	}
	text := f.MockText
	if text == "" {
		text = "Hello, how much is a taxi from Noi Bai to Hoan Kiem?"
	}
	return Translation{Text: text}, nil
}

type FakeVision struct {
	ShouldFail      bool
	MockDescription string
}

func (f *FakeVision) Observe(ctx context.Context, input VisionInput) (VisionObservation, error) {
	if f.ShouldFail {
		return VisionObservation{}, errors.New("fake vision error")
	}
	desc := f.MockDescription
	if desc == "" {
		desc = "Hóa đơn taxi 380,000 VND"
	}
	return VisionObservation{Description: desc}, nil
}

type FakeExtraction struct {
	ShouldFail bool
	MockFacts  []string
}

func (f *FakeExtraction) Extract(ctx context.Context, input ExtractionInput) (StructuredFacts, error) {
	if f.ShouldFail {
		return StructuredFacts{}, errors.New("fake extraction error")
	}
	facts := f.MockFacts
	if len(facts) == 0 {
		facts = []string{"taxi", "price_report"}
	}
	return StructuredFacts{Facts: facts}, nil
}
