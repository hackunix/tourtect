package places

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/platform/httpserver"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request, params openapi.ListPlacesParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}

	var cursor *string
	if params.Cursor != nil {
		cursor = params.Cursor
	}

	listParams := ListParams{
		RegionID: params.RegionId,
		Category: params.Category,
		Search:   params.Q,
		Cursor:   cursor,
		Limit:    limit,
		Lat:      params.Lat,
		Lon:      params.Lon,
		RadiusM:  params.RadiusM,
	}

	domainPlaces, nextCursor, err := h.service.List(ctx, listParams)
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "Invalid request", err.Error(), r.URL.Path, reqID)
		return
	}

	summaries := make([]openapi.PlaceSummary, len(domainPlaces))
	for i, dp := range domainPlaces {
		aliases := dp.Aliases
		postCount := dp.PostCount
		rating := dp.AverageRating
		fresh := dp.Freshness
		address := dp.Address
		region := dp.RegionID

		summaries[i] = openapi.PlaceSummary{
			PlaceId: dp.PlaceID,
			Name:    dp.Name,
			Category: dp.Category,
			RegionId: &region,
			Address: &address,
			Coordinates: openapi.Coordinates{
				Latitude:  dp.Latitude,
				Longitude: dp.Longitude,
			},
			Aliases: &aliases,
			PostCount: &postCount,
			AverageRating: &rating,
			Freshness: &fresh,
			CreatedAt: dp.CreatedAt,
			DistanceM: dp.DistanceM,
		}
	}

	hasMore := nextCursor != nil
	resp := openapi.PlaceListResponse{
		Items: summaries,
		Pagination: openapi.CursorInfo{
			NextCursor: nextCursor,
			HasMore:    &hasMore,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, placeId openapi.PlaceIdParam, params openapi.GetPlaceParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	id, err := uuid.Parse(placeId.String())
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "Invalid place ID", "The provided place ID is not a valid UUID", r.URL.Path, reqID)
		return
	}

	dp, err := h.service.Get(ctx, id)
	if err != nil {
		httpserver.WriteError(w, http.StatusNotFound, "Place not found", fmt.Sprintf("No place was found with ID %s", placeId.String()), r.URL.Path, reqID)
		return
	}

	aliases := dp.Aliases
	postCount := dp.PostCount
	reviewCount := dp.ReviewCount
	rating := dp.AverageRating
	fresh := dp.Freshness
	address := dp.Address
	desc := dp.Description
	phone := dp.Phone
	website := dp.Website
	oh := dp.OpeningHours
	region := dp.RegionID

	// Create copies of primitive fields so we can take address safely
	resp := openapi.PlaceDetail{
		PlaceId: dp.PlaceID,
		Name:    dp.Name,
		Category: dp.Category,
		RegionId: &region,
		Address: &address,
		Description: &desc,
		Phone: &phone,
		Website: &website,
		OpeningHours: &oh,
		Coordinates: openapi.Coordinates{
			Latitude:  dp.Latitude,
			Longitude: dp.Longitude,
		},
		Aliases: &aliases,
		PostCount: &postCount,
		ReviewCount: &reviewCount,
		AverageRating: &rating,
		Freshness: &fresh,
		CreatedAt: dp.CreatedAt,
		UpdatedAt: &dp.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
