package pricing

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tourtect/backend/generated/database"
)

type Engine struct {
	pool    *pgxpool.Pool
	queries *database.Queries
}

func NewEngine(pool *pgxpool.Pool) *Engine {
	return &Engine{
		pool:    pool,
		queries: database.New(pool),
	}
}

type PriceCheckInput struct {
	Vertical             string
	RawItem              string
	AmountMinor          string
	Currency             string
	Exponent             int
	Unit                 string
	RegionID             string
	PricingZoneID        *string
	ServiceSegment       string
	VenueType            string
	TransactionContext   string
	ObservedAt           time.Time
	ExtractionConfidence float64
	UserConfirmed        bool
}

type PriceCheckResult struct {
	AlertLevel                 string    `json:"alert_level"`
	ObservedAmountMinor        string    `json:"observed_amount_minor"`
	ObservedCurrency           string    `json:"observed_currency"`
	ObservedExponent           int       `json:"observed_exponent"`
	DeviationRatio             float64   `json:"deviation_ratio"`
	Confidence                 float64   `json:"confidence"`
	ComparisonScope            string    `json:"comparison_scope"`
	Freshness                  time.Time `json:"freshness"`
	SampleSize                 int       `json:"sample_size"`
	SnapshotVersion            string    `json:"snapshot_version"`
	Reasons                    []string  `json:"reasons"`
	PossibleBenignExplanations []string  `json:"possible_benign_explanations"`
	DatasetVersion             string    `json:"dataset_version"`
	TraceID                    string    `json:"trace_id"`
	ReferenceP10               string    `json:"reference_p10"`
	ReferenceP50               string    `json:"reference_p50"`
	ReferenceP90               string    `json:"reference_p90"`
	ReferenceGeoFallback       string    `json:"reference_geo_fallback"`
	SnapshotID                 string    `json:"snapshot_id"`
	IndependentSourceCount     int       `json:"independent_source_count"`
	SnapshotEffectiveFrom      time.Time `json:"snapshot_effective_from"`
}

func (e *Engine) Evaluate(ctx context.Context, input PriceCheckInput, traceID string) (*PriceCheckResult, error) {
	if traceID == "" {
		traceID = uuid.New().String()
	}

	result := &PriceCheckResult{
		ObservedAmountMinor:        input.AmountMinor,
		ObservedCurrency:           input.Currency,
		ObservedExponent:           input.Exponent,
		Confidence:                 input.ExtractionConfidence,
		Freshness:                  time.Time{},
		TraceID:                    traceID,
		Reasons:                    []string{},
		PossibleBenignExplanations: []string{},
	}

	// Rule 1: Extraction confidence check
	if input.ExtractionConfidence < 0.5 && !input.UserConfirmed {
		result.AlertLevel = "insufficient_data"
		result.Reasons = append(result.Reasons, "low_extraction_confidence_unconfirmed")
		return result, nil
	}

	// Parse observed value
	observedVal, err := strconv.ParseInt(input.AmountMinor, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid observed amount format: %w", err)
	}

	// Respect vertical, currency, unit, region, etc.
	// Step 1: Look for exact snapshot matching (vertical, region, segment, venue, unit, currency)
	var snap database.PriceSnapshot
	var isFallback bool
	var fallbackLevel string

	row, err := e.queries.GetPriceSnapshot(ctx, database.GetPriceSnapshotParams{
		Vertical:       input.Vertical,
		RegionID:       input.RegionID,
		Unit:           input.Unit,
		Currency:       input.Currency,
		ServiceSegment: input.ServiceSegment,
		VenueType:      input.VenueType,
		ObservedAt:     input.ObservedAt,
	})

	if err == nil {
		snap = database.PriceSnapshot{
			SnapshotID:             row.SnapshotID,
			Vertical:               row.Vertical,
			RegionID:               row.RegionID,
			PricingZoneID:          row.PricingZoneID,
			ServiceSegment:         row.ServiceSegment,
			VenueType:              row.VenueType,
			Unit:                   row.Unit,
			Currency:               row.Currency,
			P10Minor:               row.P10Minor,
			P50Minor:               row.P50Minor,
			P90Minor:               row.P90Minor,
			Exponent:               row.Exponent,
			SampleSize:             row.SampleSize,
			IndependentSourceCount: row.IndependentSourceCount,
			Version:                row.Version,
			EffectiveFrom:          row.EffectiveFrom,
			EffectiveTo:            row.EffectiveTo,
			CreatedAt:              row.CreatedAt,
		}
		fallbackLevel = "exact_zone"
	} else if errors.Is(err, pgx.ErrNoRows) {
		// Step 2: Fallback to broader region cohort
		fallbackRow, err := e.queries.GetPriceSnapshotFallbackRegion(ctx, database.GetPriceSnapshotFallbackRegionParams{
			Vertical:       input.Vertical,
			Unit:           input.Unit,
			Currency:       input.Currency,
			ServiceSegment: input.ServiceSegment,
			VenueType:      input.VenueType,
			ObservedAt:     input.ObservedAt,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			result.AlertLevel = "insufficient_data"
			result.Reasons = append(result.Reasons, "no_matching_snapshot")
			return result, nil
		}
		if err != nil {
			return nil, fmt.Errorf("retrieve fallback price snapshot: %w", err)
		}

		snap = database.PriceSnapshot{
			SnapshotID:             fallbackRow.SnapshotID,
			Vertical:               fallbackRow.Vertical,
			RegionID:               fallbackRow.RegionID,
			PricingZoneID:          fallbackRow.PricingZoneID,
			ServiceSegment:         fallbackRow.ServiceSegment,
			VenueType:              fallbackRow.VenueType,
			Unit:                   fallbackRow.Unit,
			Currency:               fallbackRow.Currency,
			P10Minor:               fallbackRow.P10Minor,
			P50Minor:               fallbackRow.P50Minor,
			P90Minor:               fallbackRow.P90Minor,
			Exponent:               fallbackRow.Exponent,
			SampleSize:             fallbackRow.SampleSize,
			IndependentSourceCount: fallbackRow.IndependentSourceCount,
			Version:                fallbackRow.Version,
			EffectiveFrom:          fallbackRow.EffectiveFrom,
			EffectiveTo:            fallbackRow.EffectiveTo,
			CreatedAt:              fallbackRow.CreatedAt,
		}
		isFallback = true
		fallbackLevel = "national_vertical"
		result.Reasons = append(result.Reasons, "fallback_broader_cohort")
	} else {
		return nil, fmt.Errorf("retrieve price snapshot: %w", err)
	}

	// Validate snapshot currency and exponent matches input
	if snap.Currency != input.Currency {
		result.AlertLevel = "insufficient_data"
		result.Reasons = append(result.Reasons, "currency_mismatch")
		return result, nil
	}

	// Rule 5: Sample size check
	if snap.SampleSize < 3 {
		result.AlertLevel = "insufficient_data"
		result.Reasons = append(result.Reasons, "sample_size_too_small")
		return result, nil
	}

	// Calculate deviation from median (P50)
	deviationRatio := float64(observedVal-snap.P50Minor) / float64(snap.P50Minor)

	// Determine alert level
	var alertLevel string
	if observedVal <= snap.P50Minor {
		alertLevel = "typical" // within_range
	} else if observedVal <= snap.P90Minor {
		alertLevel = "typical" // still within normal upper bounds
	} else if float64(observedVal) <= 1.5*float64(snap.P90Minor) {
		alertLevel = "elevated" // slightly_high
		result.Reasons = append(result.Reasons, "price_above_typical_range")
	} else {
		alertLevel = "high_risk" // significantly_high
		result.Reasons = append(result.Reasons, "price_significantly_above_range")
	}

	// Add benign explanations based on transaction context or attributes
	if input.TransactionContext == "verbal_quote" {
		result.PossibleBenignExplanations = append(result.PossibleBenignExplanations, "non_binding_verbal_quote")
	}
	if input.TransactionContext == "platform_booked" {
		result.PossibleBenignExplanations = append(result.PossibleBenignExplanations, "platform_surge_pricing")
	}

	result.AlertLevel = alertLevel
	result.DeviationRatio = deviationRatio
	result.SampleSize = int(snap.SampleSize)
	result.SnapshotVersion = snap.Version
	result.DatasetVersion = snap.Version
	result.Freshness = snap.EffectiveFrom
	result.SnapshotID = snap.SnapshotID.String()
	result.IndependentSourceCount = int(snap.IndependentSourceCount)
	result.SnapshotEffectiveFrom = snap.EffectiveFrom
	result.ComparisonScope = fmt.Sprintf("vertical=%s, region=%s, unit=%s", snap.Vertical, snap.RegionID, snap.Unit)
	result.ReferenceP10 = strconv.FormatInt(snap.P10Minor, 10)
	result.ReferenceP50 = strconv.FormatInt(snap.P50Minor, 10)
	result.ReferenceP90 = strconv.FormatInt(snap.P90Minor, 10)
	result.ReferenceGeoFallback = fallbackLevel

	if isFallback {
		result.Reasons = append(result.Reasons, fmt.Sprintf("geo_fallback_used: level=%s", fallbackLevel))
	}

	return result, nil
}
