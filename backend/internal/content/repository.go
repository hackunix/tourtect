package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	RegionID             *string
	StructuredData       json.RawMessage
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

func (r *Repository) CreateDraft(ctx context.Context, authorID uuid.UUID, postType, locale, title, body string, regionID *string, structuredData []byte, placeIDs []uuid.UUID) (*Post, error) {
	// Execute in a transaction to insert post and create links
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	var p database.Post
	err = tx.QueryRow(ctx, `INSERT INTO posts(author_id,post_type,original_locale,title,body,region_id,structured_data,moderation_status)
		VALUES($1,$2,$3,$4,$5,$6,$7,'draft') RETURNING post_id,author_id,post_type,original_locale,title,body,evidence_level,commercial_disclosure,moderation_status,created_at,updated_at,region_id,structured_data`,
		authorID, postType, locale, title, body, regionID, structuredData).Scan(&p.PostID, &p.AuthorID, &p.PostType, &p.OriginalLocale, &p.Title, &p.Body, &p.EvidenceLevel, &p.CommercialDisclosure, &p.ModerationStatus, &p.CreatedAt, &p.UpdatedAt, &p.RegionID, &p.StructuredData)
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

	if err := insertStructuredDraft(ctx, tx, p.PostID, postType, placeIDs, structuredData); err != nil {
		return nil, err
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
		RegionID:             regionID,
		StructuredData:       append(json.RawMessage(nil), structuredData...),
		EvidenceLevel:        p.EvidenceLevel,
		CommercialDisclosure: p.CommercialDisclosure,
		ModerationStatus:     p.ModerationStatus,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
		PlaceIds:             placeIDs,
	}, nil
}

func insertStructuredDraft(ctx context.Context, tx pgx.Tx, postID uuid.UUID, postType string, placeIDs []uuid.UUID, raw []byte) error {
	if postType != "review" && postType != "price_report" && postType != "scam_report" {
		return nil
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return errors.New("structured_data must be valid JSON")
	}
	firstPlace := uuid.Nil
	if len(placeIDs) > 0 {
		firstPlace = placeIDs[0]
	}
	switch postType {
	case "review":
		rating, ok := numberValue(data["overall_rating"])
		if firstPlace == uuid.Nil || !ok || rating < 1 || rating > 5 {
			return errors.New("review requires a place and overall_rating from 1 to 5")
		}
		_, err := tx.Exec(ctx, `INSERT INTO reviews(post_id,place_id,visited_at,overall_rating,price_transparency_rating,service_rating,safety_rating,value_rating)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, postID, firstPlace, nullableStringValue(data["visited_at"]), int(rating), nullableInt(data["price_transparency_rating"]), nullableInt(data["service_rating"]), nullableInt(data["safety_rating"]), nullableInt(data["value_rating"]))
		if err != nil {
			return fmt.Errorf("create review data: %w", err)
		}
	case "price_report":
		amount, ok := numberValue(data["amount_minor"])
		item, currency, unit := stringValue(data["item"]), strings.ToUpper(stringValue(data["currency"])), stringValue(data["unit"])
		observed, err := time.Parse(time.RFC3339, stringValue(data["observed_at"]))
		if !ok || amount < 0 || item == "" || len(currency) != 3 || unit == "" || err != nil {
			return errors.New("price report requires item, non-negative amount_minor, 3-letter currency, unit, and observed_at")
		}
		_, err = tx.Exec(ctx, `INSERT INTO community_price_reports(post_id,item,amount_minor,currency,unit,observed_at) VALUES($1,$2,$3,$4,$5,$6)`, postID, item, int64(amount), currency, unit, observed)
		if err != nil {
			return fmt.Errorf("create price report data: %w", err)
		}
	case "scam_report":
		state := stringValue(data["current_safety_state"])
		observed, err := time.Parse(time.RFC3339, stringValue(data["observed_at"]))
		if state == "" || err != nil {
			return errors.New("scam report requires current_safety_state and observed_at")
		}
		_, err = tx.Exec(ctx, `INSERT INTO community_scam_reports(post_id,observed_at,current_safety_state) VALUES($1,$2,$3)`, postID, observed, state)
		if err != nil {
			return fmt.Errorf("create scam report data: %w", err)
		}
	}
	return nil
}

func stringValue(value any) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}
func nullableStringValue(value any) any {
	if text := stringValue(value); text != "" {
		return text
	}
	return nil
}
func numberValue(value any) (float64, bool) { number, ok := value.(float64); return number, ok }
func nullableInt(value any) any {
	if number, ok := numberValue(value); ok {
		return int(number)
	}
	return nil
}

func (r *Repository) Publish(ctx context.Context, id, authorID uuid.UUID) (*Post, error) {
	var status string
	err := r.pool.QueryRow(ctx, `UPDATE posts
		SET moderation_status = CASE WHEN post_type = 'scam_report' THEN 'pending' ELSE 'published' END, updated_at = now()
		WHERE post_id = $1 AND author_id = $2 AND moderation_status = 'draft'
		RETURNING moderation_status`, id, authorID).Scan(&status)
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
