package content

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/platform/httpserver"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request, params openapi.ListPostsParams) {
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

	var placeID *uuid.UUID
	if params.PlaceId != nil {
		pID := uuid.UUID(*params.PlaceId)
		placeID = &pID
	}

	var postType *string
	if params.PostType != nil {
		pStr := string(*params.PostType)
		postType = &pStr
	}

	domainPosts, nextCursor, err := h.service.List(ctx, ListPostsParams{
		PlaceID:  placeID,
		PostType: postType,
		Cursor:   cursor,
		Limit:    limit,
	})
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "Invalid request", err.Error(), r.URL.Path, reqID)
		return
	}

	posts := make([]openapi.Post, len(domainPosts))
	for i, dp := range domainPosts {
		pids := make([]openapi_types.UUID, len(dp.PlaceIds))
		for idx, pid := range dp.PlaceIds {
			pids[idx] = openapi_types.UUID(pid)
		}

		posts[i] = openapi.Post{
			PostId:               dp.PostID,
			AuthorId:             dp.AuthorID,
			PostType:             openapi.PostType(dp.PostType),
			OriginalLocale:       openapi.Locale(dp.OriginalLocale),
			Title:                dp.Title,
			Body:                 dp.Body,
			EvidenceLevel:        openapi.EvidenceLevel(dp.EvidenceLevel),
			CommercialDisclosure: openapi.CommercialDisclosure(dp.CommercialDisclosure),
			ModerationStatus:     openapi.ModerationStatus(dp.ModerationStatus),
			CreatedAt:            dp.CreatedAt,
			UpdatedAt:            dp.UpdatedAt,
			PlaceIds:             &pids,
		}
	}

	hasMore := nextCursor != nil
	resp := openapi.PostListResponse{
		Items: posts,
		Pagination: openapi.CursorInfo{
			NextCursor: nextCursor,
			HasMore:    &hasMore,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CreateDraft(w http.ResponseWriter, r *http.Request, params openapi.CreateDraftParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	// Retrieve user ID from auth context
	userIDStr := httpserver.GetUserID(ctx)
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpserver.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User identity is missing or invalid", r.URL.Path, reqID)
		return
	}

	var req openapi.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Unprocessable Entity", "Invalid request body format", r.URL.Path, reqID)
		return
	}

	placeIDs := make([]uuid.UUID, 0)
	if req.PlaceIds != nil {
		for _, pid := range *req.PlaceIds {
			placeIDs = append(placeIDs, uuid.UUID(pid))
		}
	}
	structuredData := json.RawMessage(`{}`)
	if req.StructuredData != nil {
		encoded, marshalErr := json.Marshal(*req.StructuredData)
		if marshalErr != nil {
			httpserver.WriteError(w, http.StatusUnprocessableEntity, "Unprocessable Entity", "Invalid structured_data", r.URL.Path, reqID)
			return
		}
		structuredData = encoded
	}

	dp, err := h.service.CreateDraft(ctx, CreateDraftParams{
		AuthorID:       authorID,
		PostType:       string(req.PostType),
		Locale:         string(req.OriginalLocale),
		Title:          req.Title,
		Body:           req.Body,
		RegionID:       req.RegionId,
		StructuredData: structuredData,
		PlaceIds:       placeIDs,
	})
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path, reqID)
		return
	}

	pids := make([]openapi_types.UUID, len(dp.PlaceIds))
	for idx, pid := range dp.PlaceIds {
		pids[idx] = openapi_types.UUID(pid)
	}

	resp := openapi.Post{
		PostId:               dp.PostID,
		AuthorId:             dp.AuthorID,
		PostType:             openapi.PostType(dp.PostType),
		OriginalLocale:       openapi.Locale(dp.OriginalLocale),
		Title:                dp.Title,
		Body:                 dp.Body,
		EvidenceLevel:        openapi.EvidenceLevel(dp.EvidenceLevel),
		CommercialDisclosure: openapi.CommercialDisclosure(dp.CommercialDisclosure),
		ModerationStatus:     openapi.ModerationStatus(dp.ModerationStatus),
		CreatedAt:            dp.CreatedAt,
		UpdatedAt:            dp.UpdatedAt,
		PlaceIds:             &pids,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request, postId openapi.PostIdParam, params openapi.PublishPostParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	id, err := uuid.Parse(postId.String())
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "Invalid post ID", "The provided post ID is not a valid UUID", r.URL.Path, reqID)
		return
	}

	authorID, err := uuid.Parse(httpserver.GetUserID(ctx))
	if err != nil {
		httpserver.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User identity is missing or invalid", r.URL.Path, reqID)
		return
	}

	dp, err := h.service.Publish(ctx, id, authorID)
	if err != nil {
		httpserver.WriteError(w, http.StatusConflict, "Publish conflict", err.Error(), r.URL.Path, reqID)
		return
	}

	pids := make([]openapi_types.UUID, len(dp.PlaceIds))
	for idx, pid := range dp.PlaceIds {
		pids[idx] = openapi_types.UUID(pid)
	}

	resp := openapi.Post{
		PostId:               dp.PostID,
		AuthorId:             dp.AuthorID,
		PostType:             openapi.PostType(dp.PostType),
		OriginalLocale:       openapi.Locale(dp.OriginalLocale),
		Title:                dp.Title,
		Body:                 dp.Body,
		EvidenceLevel:        openapi.EvidenceLevel(dp.EvidenceLevel),
		CommercialDisclosure: openapi.CommercialDisclosure(dp.CommercialDisclosure),
		ModerationStatus:     openapi.ModerationStatus(dp.ModerationStatus),
		CreatedAt:            dp.CreatedAt,
		UpdatedAt:            dp.UpdatedAt,
		PlaceIds:             &pids,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
