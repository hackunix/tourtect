package content

import (
	"context"
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
	AuthorID uuid.UUID
	PostType string
	Locale   string
	Title    string
	Body     string
	PlaceIds []uuid.UUID
}

func (s *Service) CreateDraft(ctx context.Context, params CreateDraftParams) (*Post, error) {
	if params.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if params.Body == "" {
		return nil, errors.New("body cannot be empty")
	}

	return s.repo.CreateDraft(ctx, params.AuthorID, params.PostType, params.Locale, params.Title, params.Body, params.PlaceIds)
}

func (s *Service) Publish(ctx context.Context, id uuid.UUID) (*Post, error) {
	return s.repo.Publish(ctx, id)
}
