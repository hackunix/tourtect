package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tourtect/backend/generated/database"
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *database.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: database.New(pool),
	}
}

type Post struct {
	PostID               uuid.UUID
	AuthorID             uuid.UUID
	PostType             string
	OriginalLocale       string
	Title                string
	Body                 string
	EvidenceLevel        string
	CommercialDisclosure string
	ModerationStatus     string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	PlaceIds             []uuid.UUID
}

func (r *Repository) List(ctx context.Context, placeID *uuid.UUID, postType *string, cursor *uuid.UUID, limit int) ([]Post, error) {
	var dbPlaceID pgtype.UUID
	if placeID != nil {
		dbPlaceID = pgtype.UUID{Bytes: *placeID, Valid: true}
	}

	var cursorID pgtype.UUID
	if cursor != nil {
		cursorID = pgtype.UUID{Bytes: *cursor, Valid: true}
	}

	rows, err := r.queries.ListPublishedPosts(ctx, database.ListPublishedPostsParams{
		PlaceID:   dbPlaceID,
		PostType:  postType,
		CursorID:  cursorID,
		PageLimit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("list posts db error: %w", err)
	}

	posts := make([]Post, 0, len(rows))
	for _, row := range rows {
		posts = append(posts, Post{
			PostID:               row.PostID,
			AuthorID:             row.AuthorID,
			PostType:             row.PostType,
			OriginalLocale:       row.OriginalLocale,
			Title:                row.Title,
			Body:                 row.Body,
			EvidenceLevel:        row.EvidenceLevel,
			CommercialDisclosure: row.CommercialDisclosure,
			ModerationStatus:     row.ModerationStatus,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
			PlaceIds:             getUUIDSlice(row.PlaceIds),
		})
	}
	return posts, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Post, error) {
	row, err := r.queries.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get post db error: %w", err)
	}

	return &Post{
		PostID:               row.PostID,
		AuthorID:             row.AuthorID,
		PostType:             row.PostType,
		OriginalLocale:       row.OriginalLocale,
		Title:                row.Title,
		Body:                 row.Body,
		EvidenceLevel:        row.EvidenceLevel,
		CommercialDisclosure: row.CommercialDisclosure,
		ModerationStatus:     row.ModerationStatus,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		PlaceIds:             getUUIDSlice(row.PlaceIds),
	}, nil
}

func (r *Repository) CreateDraft(ctx context.Context, authorID uuid.UUID, postType, locale, title, body string, placeIDs []uuid.UUID) (*Post, error) {
	// Execute in a transaction to insert post and create links
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	p, err := qtx.CreateDraftPost(ctx, database.CreateDraftPostParams{
		AuthorID:       authorID,
		PostType:       postType,
		OriginalLocale: locale,
		Title:          title,
		Body:           body,
	})
	if err != nil {
		return nil, fmt.Errorf("create draft post db error: %w", err)
	}

	for _, pid := range placeIDs {
		err = qtx.LinkPostToPlace(ctx, database.LinkPostToPlaceParams{
			PostID:  p.PostID,
			PlaceID: pid,
		})
		if err != nil {
			return nil, fmt.Errorf("link post to place db error: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &Post{
		PostID:               p.PostID,
		AuthorID:             p.AuthorID,
		PostType:             p.PostType,
		OriginalLocale:       p.OriginalLocale,
		Title:                p.Title,
		Body:                 p.Body,
		EvidenceLevel:        p.EvidenceLevel,
		CommercialDisclosure: p.CommercialDisclosure,
		ModerationStatus:     p.ModerationStatus,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
		PlaceIds:             placeIDs,
	}, nil
}

func (r *Repository) Publish(ctx context.Context, id uuid.UUID) (*Post, error) {
	_, err := r.queries.PublishPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("post not found or not in draft state")
		}
		return nil, fmt.Errorf("publish post db error: %w", err)
	}

	// Fetch linked place IDs
	fullPost, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return fullPost, nil
}

// Helper function to extract UUID slices
func getUUIDSlice(i interface{}) []uuid.UUID {
	if i == nil {
		return []uuid.UUID{}
	}
	switch v := i.(type) {
	case []uuid.UUID:
		return v
	case []interface{}:
		s := make([]uuid.UUID, 0, len(v))
		for _, val := range v {
			if uStr, ok := val.(string); ok {
				if u, err := uuid.Parse(uStr); err == nil {
					s = append(s, u)
				}
			} else if uBytes, ok := val.([16]byte); ok {
				s = append(s, uuid.UUID(uBytes))
			}
		}
		return s
	}
	return []uuid.UUID{}
}
