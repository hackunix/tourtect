package pricing

import (
	"encoding/json"
	"net/http"

	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/platform/httpserver"
)

type Handler struct {
	engine *Engine
}

func NewHandler(engine *Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) CreatePriceCheck(w http.ResponseWriter, r *http.Request, params openapi.CreatePriceCheckParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	var req openapi.PriceCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Unprocessable Entity", "Invalid request body format", r.URL.Path, reqID)
		return
	}

	confidence := 1.0
	if req.ExtractionConfidence != nil {
		confidence = *req.ExtractionConfidence
	}

	confirmed := false
	if req.UserConfirmed != nil {
		confirmed = *req.UserConfirmed
	}

	input := PriceCheckInput{
		Vertical:             string(req.Vertical),
		RawItem:              req.RawItem,
		AmountMinor:          req.Money.AmountMinor,
		Currency:             req.Money.Currency,
		Exponent:             req.Money.Exponent,
		Unit:                 req.Unit,
		RegionID:             req.RegionId,
		PricingZoneID:        req.PricingZoneId,
		ServiceSegment:       string(req.ServiceSegment),
		VenueType:            string(req.VenueType),
		TransactionContext:   string(req.TransactionContext),
		ObservedAt:           req.ObservedAt,
		ExtractionConfidence: confidence,
		UserConfirmed:        confirmed,
	}

	res, err := h.engine.Evaluate(ctx, input, reqID)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path, reqID)
		return
	}

	// Format response to openapi.PriceInsight
	resp := openapi.PriceInsight{
		AlertLevel: openapi.AlertLevel(res.AlertLevel),
		Observed: openapi.Money{
			AmountMinor: res.ObservedAmountMinor,
			Currency:    res.ObservedCurrency,
			Exponent:    res.ObservedExponent,
		},
		DeviationRatio:  &res.DeviationRatio,
		Confidence:      res.Confidence,
		ComparisonScope: &res.ComparisonScope,
		Freshness:       res.Freshness,
		SampleSize:      &res.SampleSize,
		SnapshotVersion: &res.SnapshotVersion,
		Reasons:         res.Reasons,
		PossibleBenignExplanations: &res.PossibleBenignExplanations,
		DatasetVersion:  res.DatasetVersion,
		TraceId:         res.TraceID,
	}

	// Add reference if data is typical/elevated/high_risk (not insufficient_data)
	if res.AlertLevel != "insufficient_data" {
		fallbackLevel := openapi.GeoFallbackLevel(res.ReferenceGeoFallback)
		resp.Reference = &openapi.PriceReference{
			P10Minor:               &res.ReferenceP10,
			P50Minor:               &res.ReferenceP50,
			P90Minor:               &res.ReferenceP90,
			Currency:               &res.ObservedCurrency,
			Exponent:               &res.ObservedExponent,
			Unit:                   &input.Unit,
			RegionId:               &input.RegionID,
			PricingZoneId:          input.PricingZoneID,
			ServiceSegment:         (*openapi.ServiceSegment)(&input.ServiceSegment),
			VenueType:              (*openapi.VenueType)(&input.VenueType),
			GeoFallbackLevel:       &fallbackLevel,
			EffectiveSampleSize:    &res.SampleSize,
			IndependentSourceCount: &res.SampleSize, // mapped for slice completeness
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
