package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type ListPostsParams struct {
	PlaceID  *uuid.UUID
	PostType *string
	Cursor   *string
	Limit    int
}

func (s *Service) List(ctx context.Context, params ListPostsParams) ([]Post, *string, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	var cursorUUID *uuid.UUID
	if params.Cursor != nil && *params.Cursor != "" {
		id, err := uuid.Parse(*params.Cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid cursor: %w", err)
		}
		cursorUUID = &id
	}

	posts, err := s.repo.List(ctx, params.PlaceID, params.PostType, cursorUUID, params.Limit)
	if err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if len(posts) == params.Limit {
		lastID := posts[len(posts)-1].PostID.String()
		nextCursor = &lastID
	}

	return posts, nextCursor, nil
}

type CreateDraftParams struct {
	AuthorID       uuid.UUID
	PostType       string
	Locale         string
	Title          string
	Body           string
	RegionID       *string
	StructuredData json.RawMessage
	PlaceIds       []uuid.UUID
}

func (s *Service) CreateDraft(ctx context.Context, params CreateDraftParams) (*Post, error) {
	if params.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if params.Body == "" {
		return nil, errors.New("body cannot be empty")
	}
	allowed := map[string]bool{"discussion": true, "question": true, "review": true, "price_report": true, "scam_report": true, "tip": true}
	if !allowed[params.PostType] {
		return nil, errors.New("post type is not available in the traveler composer")
	}
	if (params.PostType == "review" || params.PostType == "price_report" || params.PostType == "scam_report") && len(params.StructuredData) <= 2 {
		return nil, errors.New("structured_data is required for this post type")
	}

	return s.repo.CreateDraft(ctx, params.AuthorID, params.PostType, params.Locale, params.Title, params.Body, params.RegionID, params.StructuredData, params.PlaceIds)
}

func (s *Service) Publish(ctx context.Context, id, authorID uuid.UUID) (*Post, error) {
	return s.repo.Publish(ctx, id, authorID)
}
