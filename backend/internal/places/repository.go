package places

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

type Place struct {
	PlaceID       uuid.UUID
	Name          string
	Category      string
	RegionID      string
	Address       string
	Description   string
	Phone         string
	Website       string
	OpeningHours  string
	Latitude      float64
	Longitude     float64
	PostCount     int
	ReviewCount   int
	AverageRating float64
	Freshness     time.Time
	CreatedAt     time.Time
	Aliases       []string
	DistanceM     *float64
}

func (r *Repository) List(ctx context.Context, regionID, category, search *string, cursor *uuid.UUID, limit int) ([]Place, error) {
	var cursorID pgtype.UUID
	if cursor != nil {
		cursorID = pgtype.UUID{Bytes: *cursor, Valid: true}
	}

	rows, err := r.queries.ListPlaces(ctx, database.ListPlacesParams{
		RegionID:    regionID,
		Category:    category,
		SearchQuery: search,
		CursorID:    cursorID,
		PageLimit:   int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("list places db error: %w", err)
	}

	places := make([]Place, 0, len(rows))
	for _, row := range rows {
		places = append(places, Place{
			PlaceID:       row.PlaceID,
			Name:          row.Name,
			Category:      row.Category,
			RegionID:      row.RegionID,
			Address:       getString(row.Address),
			Latitude:      getFloat64(row.Latitude),
			Longitude:     getFloat64(row.Longitude),
			PostCount:     int(row.PostCount),
			AverageRating: getNumeric(row.AverageRating),
			Freshness:     getTime(row.Freshness),
			CreatedAt:     row.CreatedAt,
			Aliases:       getStringSlice(row.Aliases),
		})
	}
	return places, nil
}

func (r *Repository) ListNearby(ctx context.Context, lat, lon float64, radiusM int, limit int) ([]Place, error) {
	rows, err := r.queries.ListPlacesNearby(ctx, database.ListPlacesNearbyParams{
		Lat:       lat,
		Lon:       lon,
		RadiusM:   float64(radiusM),
		PageLimit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("list places nearby db error: %w", err)
	}

	places := make([]Place, 0, len(rows))
	for _, row := range rows {
		dist := getFloat64(row.DistanceM)
		places = append(places, Place{
			PlaceID:       row.PlaceID,
			Name:          row.Name,
			Category:      row.Category,
			RegionID:      row.RegionID,
			Address:       getString(row.Address),
			Latitude:      getFloat64(row.Latitude),
			Longitude:     getFloat64(row.Longitude),
			PostCount:     int(row.PostCount),
			AverageRating: getNumeric(row.AverageRating),
			Freshness:     getTime(row.Freshness),
			CreatedAt:     row.CreatedAt,
			Aliases:       getStringSlice(row.Aliases),
			DistanceM:     &dist,
		})
	}
	return places, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Place, error) {
	row, err := r.queries.GetPlace(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get place db error: %w", err)
	}

	return &Place{
		PlaceID:       row.PlaceID,
		Name:          row.Name,
		Category:      row.Category,
		RegionID:      row.RegionID,
		Address:       getString(row.Address),
		Description:   getString(row.Description),
		Phone:         getString(row.Phone),
		Website:       getString(row.Website),
		OpeningHours:  getString(row.OpeningHours),
		Latitude:      getFloat64(row.Latitude),
		Longitude:     getFloat64(row.Longitude),
		PostCount:     int(row.PostCount),
		ReviewCount:   int(row.ReviewCount),
		AverageRating: getNumeric(row.AverageRating),
		Freshness:     getTime(row.Freshness),
		CreatedAt:     row.CreatedAt,
		Aliases:       getStringSlice(row.Aliases),
	}, nil
}

// Helper parsing functions

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getFloat64(i interface{}) float64 {
	switch v := i.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int:
		return float64(v)
	}
	return 0.0
}

func getNumeric(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	f, err := n.Float64Value()
	if err != nil {
		return 0.0
	}
	return f.Float64
}

func getTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func getStringSlice(i interface{}) []string {
	if i == nil {
		return []string{}
	}
	// Convert array/slice representation if needed
	switch v := i.(type) {
	case []string:
		return v
	case []interface{}:
		s := make([]string, len(v))
		for idx, val := range v {
			s[idx] = fmt.Sprintf("%v", val)
		}
		return s
	}
	return []string{}
}
