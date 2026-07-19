package retrieval

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/tourtect/backend/internal/content"
	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/places"
)

const (
	MaxSources           = 6
	MaxContextCharacters = 6000
)

type Service struct {
	places *places.Service
	posts  *content.Service
	now    func() time.Time
}

func NewService(placeService *places.Service, posts *content.Service) *Service {
	return &Service{places: placeService, posts: posts, now: time.Now}
}

type Result struct {
	Place    *places.Place
	Evidence []model.Evidence
}

func (s *Service) Retrieve(ctx context.Context, text, locale, explicitPlaceID string) (Result, error) {
	var place *places.Place
	if explicitPlaceID != "" {
		if id, err := uuid.Parse(explicitPlaceID); err == nil {
			place, _ = s.places.Get(ctx, id)
		}
	}
	if place == nil {
		for _, query := range placeQueries(text) {
			found, _, err := s.places.List(ctx, places.ListParams{Search: &query, Limit: 3})
			if err != nil {
				return Result{}, err
			}
			if len(found) > 0 {
				selected := found[0]
				place = &selected
				break
			}
		}
	}
	result := Result{Place: place, Evidence: []model.Evidence{}}
	if place == nil {
		return result, nil
	}
	result.Evidence = append(result.Evidence, model.Evidence{ID: uuid.NewString(), SourceType: "place_record", SourceID: place.PlaceID.String(), Title: place.Name, Summary: bounded(redactPII(place.Address), 320), ObservedAt: &place.Freshness, Freshness: freshness(place.Freshness, s.now()), EvidenceLevel: "verified", SourceURL: "/places/" + place.PlaceID.String()})
	posts, _, err := s.posts.List(ctx, content.ListPostsParams{PlaceID: &place.PlaceID, Limit: MaxSources})
	if err != nil {
		return Result{}, err
	}
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].OriginalLocale == locale && posts[j].OriginalLocale != locale || posts[i].UpdatedAt.After(posts[j].UpdatedAt)
	})
	seen, used := map[string]bool{}, utf8.RuneCountInString(result.Evidence[0].Summary)
	for _, post := range posts {
		if len(result.Evidence) >= MaxSources || seen[post.PostID.String()] {
			continue
		}
		summary := bounded(redactPII(post.Body), 700)
		if used+utf8.RuneCountInString(summary) > MaxContextCharacters {
			break
		}
		seen[post.PostID.String()] = true
		used += utf8.RuneCountInString(summary)
		level := "community"
		if post.EvidenceLevel == "verified_source" || post.EvidenceLevel == "verified_receipt" {
			level = "verified"
		}
		observed := post.UpdatedAt
		result.Evidence = append(result.Evidence, model.Evidence{ID: uuid.NewString(), SourceType: "community_post", SourceID: post.PostID.String(), Title: bounded(post.Title, 180), Summary: summary, ObservedAt: &observed, Freshness: freshness(observed, s.now()), EvidenceLevel: level, SourceURL: "/posts/" + post.PostID.String()})
	}
	return result, nil
}

func placeQueries(text string) []string {
	lower := strings.ToLower(text)
	queries := []string{}
	for _, candidate := range []struct{ match, query string }{{"noi bai", "Noi Bai"}, {"nội bài", "Nội Bài"}, {"hoan kiem", "Hoan Kiem"}, {"hoàn kiếm", "Hoàn Kiếm"}, {"dong xuan", "Dong Xuan"}, {"đồng xuân", "Đồng Xuân"}} {
		if strings.Contains(lower, candidate.match) {
			queries = append(queries, candidate.query)
		}
	}
	if len(queries) == 0 && len(strings.Fields(text)) <= 8 {
		queries = append(queries, strings.TrimSpace(text))
	}
	return queries
}

var emailPattern = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
var phonePattern = regexp.MustCompile(`(?:\+?[0-9][0-9\s().-]{7,}[0-9])`)

func redactPII(value string) string {
	value = emailPattern.ReplaceAllString(value, "[redacted email]")
	return phonePattern.ReplaceAllString(value, "[redacted phone]")
}
func bounded(value string, n int) string {
	r := []rune(strings.TrimSpace(value))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n]) + "…"
}
func freshness(observed, now time.Time) string {
	if observed.IsZero() {
		return "unknown"
	}
	age := now.Sub(observed)
	if age <= 30*24*time.Hour {
		return "fresh"
	}
	if age <= 180*24*time.Hour {
		return "aging"
	}
	return "stale"
}
