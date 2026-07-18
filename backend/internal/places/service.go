package places

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

type ListParams struct {
	RegionID *string
	Category *string
	Search   *string
	Cursor   *string
	Limit    int
	Lat      *float64
	Lon      *float64
	RadiusM  *int
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Place, *string, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Geolocation distance search if lat, lon, radius are provided
	if params.Lat != nil && params.Lon != nil {
		radius := 5000
		if params.RadiusM != nil {
			radius = *params.RadiusM
		}
		places, err := s.repo.ListNearby(ctx, *params.Lat, *params.Lon, radius, params.Limit)
		return places, nil, err
	}

	// Standard paginated lookup
	var cursorUUID *uuid.UUID
	if params.Cursor != nil && *params.Cursor != "" {
		id, err := uuid.Parse(*params.Cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid cursor: %w", err)
		}
		cursorUUID = &id
	}

	places, err := s.repo.List(ctx, params.RegionID, params.Category, params.Search, cursorUUID, params.Limit)
	if err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if len(places) == params.Limit {
		// Set next cursor as the ID of the last element
		lastID := places[len(places)-1].PlaceID.String()
		nextCursor = &lastID
	}

	return places, nextCursor, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Place, error) {
	place, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if place == nil {
		return nil, errors.New("place not found")
	}
	return place, nil
}
