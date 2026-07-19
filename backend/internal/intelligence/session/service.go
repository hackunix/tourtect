package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/tourtect/backend/internal/intelligence/model"
)

const (
	DefaultTTL         = 30 * time.Minute
	MaxSerializedSize  = 256 * 1024
	MaxRecentResponses = 20
	MaxProcessedIDs    = 64
)

type Service struct {
	store Store
	ttl   time.Duration
	now   func() time.Time
	locks sync.Map
}

// Lock serializes state transitions for one session in this API process. Redis
// remains the source of truth; a future multi-replica deployment must replace
// this with versioned compare-and-set semantics.
func (s *Service) Lock(id string) func() {
	value, _ := s.locks.LoadOrStore(id, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func NewService(store Store, ttl time.Duration) *Service {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Service{store: store, ttl: ttl, now: time.Now}
}

func (s *Service) Create(ctx context.Context, userID, locale, targetLocale, placeID, region, mode string, processingConsent bool) (*model.Session, error) {
	if strings.TrimSpace(locale) == "" {
		return nil, errors.New("locale is required")
	}
	if mode == "" {
		mode = "text"
	}
	now := s.now().UTC()
	value := &model.Session{
		ID: uuid.NewString(), UserID: userID, Version: model.SessionVersion,
		CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(s.ttl),
		Context:         model.SessionContext{Locale: locale, TargetLocale: targetLocale, PlaceID: placeID, ApproximateRegion: region, InteractionMode: mode, CurrentSafetyState: "unknown", ConsentState: model.ConsentState{Processing: processingConsent}},
		RecentResponses: []model.Response{}, ProcessedMessageIDs: []string{}, Confirmations: map[string]*model.Confirmation{},
	}
	if err := s.save(ctx, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Service) GetOwned(ctx context.Context, id, userID string) (*model.Session, error) {
	value, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if value.UserID != userID || value.Version != model.SessionVersion {
		return nil, ErrNotFound
	}
	return value, nil
}

func (s *Service) DeleteOwned(ctx context.Context, id, userID string) error {
	if _, err := s.GetOwned(ctx, id, userID); err != nil {
		return err
	}
	return s.store.Delete(ctx, id)
}

func (s *Service) HasMessage(value *model.Session, id string) bool {
	return slices.Contains(value.ProcessedMessageIDs, id)
}

func (s *Service) AppendResponse(ctx context.Context, value *model.Session, messageID string, response model.Response) error {
	value.RecentResponses = append(value.RecentResponses, response)
	if len(value.RecentResponses) > MaxRecentResponses {
		value.RecentResponses = value.RecentResponses[len(value.RecentResponses)-MaxRecentResponses:]
	}
	value.ProcessedMessageIDs = append(value.ProcessedMessageIDs, messageID)
	if len(value.ProcessedMessageIDs) > MaxProcessedIDs {
		value.ProcessedMessageIDs = value.ProcessedMessageIDs[len(value.ProcessedMessageIDs)-MaxProcessedIDs:]
	}
	value.Context.CurrentSafetyState = response.SafetyState
	if response.RequestedConfirmation != nil {
		if value.Confirmations == nil {
			value.Confirmations = map[string]*model.Confirmation{}
		}
		value.Confirmations[response.RequestedConfirmation.ID] = response.RequestedConfirmation
	}
	return s.save(ctx, value)
}

func (s *Service) Save(ctx context.Context, value *model.Session) error { return s.save(ctx, value) }

func (s *Service) save(ctx context.Context, value *model.Session) error {
	value.UpdatedAt = s.now().UTC()
	value.ExpiresAt = value.UpdatedAt.Add(s.ttl)
	redact(value)
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(b) > MaxSerializedSize {
		return fmt.Errorf("assistant session exceeds %d bytes", MaxSerializedSize)
	}
	return s.store.Put(ctx, value, s.ttl)
}

func redact(value *model.Session) {
	value.Context.UserConfirmedFacts = boundedStrings(value.Context.UserConfirmedFacts, 32, 512)
	if len(value.Context.ActiveCaptureIDs) > 8 {
		value.Context.ActiveCaptureIDs = value.Context.ActiveCaptureIDs[len(value.Context.ActiveCaptureIDs)-8:]
	}
	for i := range value.RecentResponses {
		value.RecentResponses[i].Message = bounded(value.RecentResponses[i].Message, 8000)
	}
}

func boundedStrings(values []string, maxItems, maxRunes int) []string {
	if len(values) > maxItems {
		values = values[len(values)-maxItems:]
	}
	for i := range values {
		values[i] = bounded(values[i], maxRunes)
	}
	return values
}

func bounded(value string, maxRunes int) string {
	if utf8.RuneCountInString(value) <= maxRunes {
		return value
	}
	r := []rune(value)
	return string(r[:maxRunes])
}
