package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/tourtect/backend/internal/platform/httpserver"
)

type authorSummary struct {
	PrincipalID uuid.UUID `json:"principal_id"`
	DisplayName string    `json:"display_name"`
}

type placeAttachment struct {
	PlaceID   uuid.UUID `json:"place_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	RegionID  string    `json:"region_id"`
	Freshness time.Time `json:"freshness"`
}

type communityPost struct {
	PostID               uuid.UUID         `json:"post_id"`
	AuthorID             uuid.UUID         `json:"author_id"`
	Author               authorSummary     `json:"author"`
	PostType             string            `json:"post_type"`
	OriginalLocale       string            `json:"original_locale"`
	Title                string            `json:"title"`
	Body                 string            `json:"body"`
	RegionID             *string           `json:"region_id,omitempty"`
	EvidenceLevel        string            `json:"evidence_level"`
	CommercialDisclosure string            `json:"commercial_disclosure"`
	ModerationStatus     string            `json:"moderation_status"`
	StructuredData       json.RawMessage   `json:"structured_data"`
	Places               []placeAttachment `json:"places"`
	UsefulCount          int               `json:"useful_count"`
	CommentCount         int               `json:"comment_count"`
	ViewerUseful         bool              `json:"viewer_useful"`
	ViewerSaved          bool              `json:"viewer_saved"`
	ReasonCodes          []string          `json:"reason_codes,omitempty"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
}

type cursorInfo struct {
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}

func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "latest"
	}
	validModes := map[string]bool{"following": true, "nearby": true, "latest": true, "trending": true, "safety": true}
	if !validModes[mode] {
		writeCommunityError(w, r, http.StatusUnprocessableEntity, "Invalid feed mode", "mode must be following, nearby, latest, trending, or safety")
		return
	}
	regionID := strings.TrimSpace(r.URL.Query().Get("region_id"))
	if mode == "nearby" && regionID == "" {
		writeCommunityError(w, r, http.StatusUnprocessableEntity, "Region required", "Choose a region before opening the nearby feed; precise location remains opt-in")
		return
	}
	limit := parseLimit(r, 20)
	viewerID, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}
	var cursor *uuid.UUID
	if raw := r.URL.Query().Get("cursor"); raw != "" {
		parsed, parseErr := uuid.Parse(raw)
		if parseErr != nil {
			writeCommunityError(w, r, http.StatusBadRequest, "Invalid cursor", "cursor must be a post UUID")
			return
		}
		cursor = &parsed
	}

	posts, err := h.listFeed(r, viewerID, mode, regionID, cursor, limit+1)
	if err != nil {
		writeCommunityError(w, r, http.StatusInternalServerError, "Feed unavailable", err.Error())
		return
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	var next *string
	if hasMore && len(posts) > 0 {
		value := posts[len(posts)-1].PostID.String()
		next = &value
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": posts, "pagination": cursorInfo{NextCursor: next, HasMore: hasMore}})
}

func (h *Handler) listFeed(r *http.Request, viewerID uuid.UUID, mode, regionID string, cursor *uuid.UUID, limit int) ([]communityPost, error) {
	filter, order, reason := feedPolicy(mode)

	query := fmt.Sprintf(`
		SELECT p.post_id, p.author_id, pr.display_name, p.post_type, p.original_locale, p.title, p.body,
			p.region_id, p.evidence_level, p.commercial_disclosure, p.moderation_status,
			p.structured_data, p.created_at, p.updated_at,
			(SELECT count(*) FROM post_votes pv WHERE pv.post_id=p.post_id),
			(SELECT count(*) FROM post_comments pc WHERE pc.post_id=p.post_id AND pc.moderation_status='published'),
			EXISTS(SELECT 1 FROM post_votes pv WHERE pv.post_id=p.post_id AND pv.principal_id=$1),
			EXISTS(SELECT 1 FROM saved_posts sp WHERE sp.post_id=p.post_id AND sp.principal_id=$1)
		FROM posts p JOIN principals pr ON pr.principal_id=p.author_id
		WHERE p.moderation_status='published' AND ($2::text IS NULL OR TRUE) AND (%s)
			AND ($3::uuid IS NULL OR p.created_at < (SELECT created_at FROM posts WHERE post_id=$3))
			AND NOT EXISTS (SELECT 1 FROM principal_blocks pb WHERE pb.blocker_id=$1 AND pb.blocked_id=p.author_id)
		ORDER BY %s LIMIT $4`, filter, order)
	rows, err := h.service.repo.pool.Query(r.Context(), query, viewerID, nullableString(regionID), cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityPost, 0)
	for rows.Next() {
		var item communityPost
		if err := rows.Scan(&item.PostID, &item.AuthorID, &item.Author.DisplayName, &item.PostType, &item.OriginalLocale,
			&item.Title, &item.Body, &item.RegionID, &item.EvidenceLevel, &item.CommercialDisclosure,
			&item.ModerationStatus, &item.StructuredData, &item.CreatedAt, &item.UpdatedAt, &item.UsefulCount,
			&item.CommentCount, &item.ViewerUseful, &item.ViewerSaved); err != nil {
			return nil, err
		}
		item.Author.PrincipalID = item.AuthorID
		item.ReasonCodes = []string{reason}
		item.Places, err = h.postPlaces(r, item.PostID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func feedPolicy(mode string) (filter, order, reason string) {
	filter = "TRUE"
	order = "p.created_at DESC, p.post_id DESC"
	reason = "latest"
	switch mode {
	case "nearby":
		filter = "p.region_id = $2 OR EXISTS (SELECT 1 FROM post_place_links ppln JOIN places pln ON pln.place_id=ppln.place_id WHERE ppln.post_id=p.post_id AND pln.region_id=$2)"
		reason = "nearby_region"
	case "following":
		filter = `(EXISTS (SELECT 1 FROM follows f WHERE f.principal_id=$1 AND f.target_principal_id=p.author_id)
			OR EXISTS (SELECT 1 FROM follows f JOIN post_place_links pplf ON pplf.place_id=f.target_place_id WHERE f.principal_id=$1 AND pplf.post_id=p.post_id))`
		reason = "followed_source"
	case "trending":
		order = "((SELECT count(*) FROM post_votes pv WHERE pv.post_id=p.post_id) * 3 + (SELECT count(*) FROM post_comments pc WHERE pc.post_id=p.post_id) * 2) / GREATEST(1, EXTRACT(EPOCH FROM (now()-p.created_at))/3600 + 2) DESC, p.created_at DESC"
		reason = "community_usefulness"
	case "safety":
		filter = "p.post_type IN ('official_alert','scam_report')"
		order = "CASE p.evidence_level WHEN 'verified_source' THEN 0 WHEN 'verified_receipt' THEN 1 WHEN 'metadata' THEN 2 ELSE 3 END, p.created_at DESC"
		reason = "safety_priority"
	}
	return filter, order, reason
}

func (h *Handler) postPlaces(r *http.Request, postID uuid.UUID) ([]placeAttachment, error) {
	rows, err := h.service.repo.pool.Query(r.Context(), `SELECT pl.place_id, pl.name, pl.category, pl.region_id, pl.freshness
		FROM places pl JOIN post_place_links ppl ON ppl.place_id=pl.place_id WHERE ppl.post_id=$1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	places := make([]placeAttachment, 0)
	for rows.Next() {
		var place placeAttachment
		if err := rows.Scan(&place.PlaceID, &place.Name, &place.Category, &place.RegionID, &place.Freshness); err != nil {
			return nil, err
		}
		places = append(places, place)
	}
	return places, rows.Err()
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 2 {
		writeCommunityError(w, r, http.StatusUnprocessableEntity, "Search term required", "q must contain at least two characters")
		return
	}
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "top"
	}
	like := "%" + q + "%"
	result := map[string]any{"query": q, "tab": tab, "places": []any{}, "posts": []any{}}
	if tab == "top" || tab == "places" {
		rows, err := h.service.repo.pool.Query(r.Context(), `SELECT place_id, name, category, region_id, address, freshness, ST_Y(coordinates::geometry), ST_X(coordinates::geometry), created_at
			FROM places WHERE name ILIKE $1 OR address ILIKE $1 OR EXISTS (SELECT 1 FROM place_aliases pa WHERE pa.place_id=places.place_id AND pa.alias ILIKE $1)
			ORDER BY similarity(name, $2) DESC LIMIT 20`, like, q)
		if err != nil {
			writeCommunityError(w, r, http.StatusInternalServerError, "Search unavailable", err.Error())
			return
		}
		places := make([]map[string]any, 0)
		for rows.Next() {
			var id uuid.UUID
			var name, category, region string
			var address *string
			var freshness, created time.Time
			var latitude, longitude float64
			if err := rows.Scan(&id, &name, &category, &region, &address, &freshness, &latitude, &longitude, &created); err != nil {
				rows.Close()
				writeCommunityError(w, r, 500, "Search unavailable", err.Error())
				return
			}
			places = append(places, map[string]any{"place_id": id, "name": name, "category": category, "region_id": region, "address": address, "freshness": freshness, "coordinates": map[string]float64{"latitude": latitude, "longitude": longitude}, "created_at": created})
		}
		rows.Close()
		result["places"] = places
	}
	if tab != "places" {
		typeFilter := "TRUE"
		if tab == "price_reports" {
			typeFilter = "p.post_type='price_report'"
		}
		if tab == "safety" {
			typeFilter = "p.post_type IN ('scam_report','official_alert')"
		}
		rows, err := h.service.repo.pool.Query(r.Context(), fmt.Sprintf(`SELECT p.post_id, p.author_id, pr.display_name, p.post_type, p.title, p.body, p.original_locale, p.evidence_level, p.commercial_disclosure, p.moderation_status, p.created_at, p.updated_at
			FROM posts p JOIN principals pr ON pr.principal_id=p.author_id WHERE p.moderation_status='published' AND %s AND (p.title ILIKE $1 OR p.body ILIKE $1)
			ORDER BY similarity(p.title, $2) DESC, p.created_at DESC LIMIT 30`, typeFilter), like, q)
		if err != nil {
			writeCommunityError(w, r, 500, "Search unavailable", err.Error())
			return
		}
		posts := make([]map[string]any, 0)
		for rows.Next() {
			var id, authorID uuid.UUID
			var authorName, pt, title, body, locale, evidence, commercial, moderation string
			var created, updated time.Time
			if err := rows.Scan(&id, &authorID, &authorName, &pt, &title, &body, &locale, &evidence, &commercial, &moderation, &created, &updated); err != nil {
				rows.Close()
				writeCommunityError(w, r, 500, "Search unavailable", err.Error())
				return
			}
			posts = append(posts, map[string]any{"post_id": id, "author_id": authorID, "author": authorSummary{PrincipalID: authorID, DisplayName: authorName}, "post_type": pt, "title": title, "body": body, "original_locale": locale, "evidence_level": evidence, "commercial_disclosure": commercial, "moderation_status": moderation, "created_at": created, "updated_at": updated})
		}
		rows.Close()
		result["posts"] = posts
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Comments(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(r.PathValue("postId"))
	if err != nil {
		writeCommunityError(w, r, 400, "Invalid post ID", "postId must be a UUID")
		return
	}
	if r.Method == http.MethodGet {
		rows, err := h.service.repo.pool.Query(r.Context(), `SELECT c.comment_id, c.author_id, p.display_name, c.parent_comment_id, c.body, c.moderation_status, c.created_at, c.updated_at
			FROM post_comments c JOIN principals p ON p.principal_id=c.author_id WHERE c.post_id=$1 AND c.moderation_status='published' ORDER BY c.created_at LIMIT $2`, postID, parseLimit(r, 50))
		if err != nil {
			writeCommunityError(w, r, 500, "Comments unavailable", err.Error())
			return
		}
		defer rows.Close()
		items := make([]map[string]any, 0)
		for rows.Next() {
			var id, author uuid.UUID
			var name, body, status string
			var parent *uuid.UUID
			var created, updated time.Time
			if err := rows.Scan(&id, &author, &name, &parent, &body, &status, &created, &updated); err != nil {
				writeCommunityError(w, r, 500, "Comments unavailable", err.Error())
				return
			}
			items = append(items, map[string]any{"comment_id": id, "post_id": postID, "author": authorSummary{PrincipalID: author, DisplayName: name}, "parent_comment_id": parent, "body": body, "moderation_status": status, "created_at": created, "updated_at": updated})
		}
		writeJSON(w, 200, map[string]any{"items": items, "pagination": cursorInfo{HasMore: false}})
		return
	}
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	var body struct {
		Body            string     `json:"body"`
		ParentCommentID *uuid.UUID `json:"parent_comment_id"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil || strings.TrimSpace(body.Body) == "" {
		writeCommunityError(w, r, 422, "Invalid comment", "body is required")
		return
	}
	if body.ParentCommentID != nil {
		var parentPost uuid.UUID
		if err := h.service.repo.pool.QueryRow(r.Context(), `SELECT post_id FROM post_comments WHERE comment_id=$1`, body.ParentCommentID).Scan(&parentPost); err != nil || parentPost != postID {
			writeCommunityError(w, r, 422, "Invalid parent", "parent comment must belong to the same post")
			return
		}
	}
	var id uuid.UUID
	var created time.Time
	err = h.service.repo.pool.QueryRow(r.Context(), `INSERT INTO post_comments(post_id,author_id,parent_comment_id,body) VALUES($1,$2,$3,$4) RETURNING comment_id,created_at`, postID, viewer, body.ParentCommentID, strings.TrimSpace(body.Body)).Scan(&id, &created)
	if err != nil {
		writeCommunityError(w, r, 400, "Comment failed", err.Error())
		return
	}
	_, _ = h.service.repo.pool.Exec(r.Context(), `INSERT INTO notifications(principal_id,kind,actor_id,post_id,comment_id,message)
		SELECT author_id,'reply',$2,$1,$3,'Có phản hồi mới trong bài viết của bạn' FROM posts WHERE post_id=$1 AND author_id<>$2`, postID, viewer, id)
	writeJSON(w, 201, map[string]any{"comment_id": id, "post_id": postID, "author_id": viewer, "parent_comment_id": body.ParentCommentID, "body": strings.TrimSpace(body.Body), "moderation_status": "published", "created_at": created, "updated_at": created})
}

func (h *Handler) UsefulVote(w http.ResponseWriter, r *http.Request) {
	h.togglePostRelation(w, r, "post_votes", "vote")
}
func (h *Handler) SavedPost(w http.ResponseWriter, r *http.Request) {
	h.togglePostRelation(w, r, "saved_posts", "saved")
}

func (h *Handler) togglePostRelation(w http.ResponseWriter, r *http.Request, table, label string) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	postID, err := uuid.Parse(r.PathValue("postId"))
	if err != nil {
		writeCommunityError(w, r, 400, "Invalid post ID", "postId must be a UUID")
		return
	}
	if r.Method == http.MethodDelete {
		_, err = h.service.repo.pool.Exec(r.Context(), fmt.Sprintf("DELETE FROM %s WHERE principal_id=$1 AND post_id=$2", table), viewer, postID)
	} else {
		_, err = h.service.repo.pool.Exec(r.Context(), fmt.Sprintf("INSERT INTO %s(principal_id,post_id) VALUES($1,$2) ON CONFLICT DO NOTHING", table), viewer, postID)
	}
	if err != nil {
		writeCommunityError(w, r, 400, "Interaction failed", err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"post_id": postID, label: r.Method != http.MethodDelete})
}

func (h *Handler) SavedList(w http.ResponseWriter, r *http.Request) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	rows, err := h.service.repo.pool.Query(r.Context(), `SELECT p.post_id,p.post_type,p.title,p.body,p.original_locale,p.evidence_level,p.created_at,sp.created_at
		FROM saved_posts sp JOIN posts p ON p.post_id=sp.post_id WHERE sp.principal_id=$1 AND p.moderation_status='published' ORDER BY sp.created_at DESC`, viewer)
	if err != nil {
		writeCommunityError(w, r, 500, "Saved unavailable", err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id uuid.UUID
		var pt, title, body, locale, evidence string
		var created, saved time.Time
		if err := rows.Scan(&id, &pt, &title, &body, &locale, &evidence, &created, &saved); err != nil {
			writeCommunityError(w, r, 500, "Saved unavailable", err.Error())
			return
		}
		items = append(items, map[string]any{"post_id": id, "post_type": pt, "title": title, "body": body, "original_locale": locale, "evidence_level": evidence, "created_at": created, "saved_at": saved})
	}
	writeJSON(w, 200, map[string]any{"items": items, "pagination": cursorInfo{HasMore: false}})
}

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	if r.Method == http.MethodPatch {
		var body struct {
			NotificationIDs []uuid.UUID `json:"notification_ids"`
			Read            bool        `json:"read"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			writeCommunityError(w, r, 422, "Invalid notification update", "notification_ids is required")
			return
		}
		_, err = h.service.repo.pool.Exec(r.Context(), `UPDATE notifications SET read_at=CASE WHEN $3 THEN now() ELSE NULL END WHERE principal_id=$1 AND notification_id=ANY($2::uuid[])`, viewer, body.NotificationIDs, body.Read)
		if err != nil {
			writeCommunityError(w, r, 400, "Notification update failed", err.Error())
			return
		}
		writeJSON(w, 200, map[string]any{"updated": len(body.NotificationIDs)})
		return
	}
	rows, err := h.service.repo.pool.Query(r.Context(), `SELECT notification_id,kind,actor_id,post_id,comment_id,message,read_at,created_at FROM notifications WHERE principal_id=$1 ORDER BY created_at DESC LIMIT $2`, viewer, parseLimit(r, 30))
	if err != nil {
		writeCommunityError(w, r, 500, "Notifications unavailable", err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id uuid.UUID
		var kind, msg string
		var actor, post, comment *uuid.UUID
		var read *time.Time
		var created time.Time
		if err := rows.Scan(&id, &kind, &actor, &post, &comment, &msg, &read, &created); err != nil {
			writeCommunityError(w, r, 500, "Notifications unavailable", err.Error())
			return
		}
		items = append(items, map[string]any{"notification_id": id, "kind": kind, "actor_id": actor, "post_id": post, "comment_id": comment, "message": msg, "read_at": read, "created_at": created})
	}
	writeJSON(w, 200, map[string]any{"items": items, "pagination": cursorInfo{HasMore: false}})
}

func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	var body struct {
		TargetType string    `json:"target_type"`
		TargetID   uuid.UUID `json:"target_id"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil || (body.TargetType != "principal" && body.TargetType != "place") {
		writeCommunityError(w, r, 422, "Invalid follow", "target_type must be principal or place")
		return
	}
	column := "target_principal_id"
	if body.TargetType == "place" {
		column = "target_place_id"
	}
	if r.Method == http.MethodDelete {
		_, err = h.service.repo.pool.Exec(r.Context(), fmt.Sprintf("DELETE FROM follows WHERE principal_id=$1 AND %s=$2", column), viewer, body.TargetID)
	} else {
		_, err = h.service.repo.pool.Exec(r.Context(), fmt.Sprintf("INSERT INTO follows(principal_id,target_type,%s) VALUES($1,$2,$3) ON CONFLICT DO NOTHING", column), viewer, body.TargetType, body.TargetID)
	}
	if err != nil {
		writeCommunityError(w, r, 400, "Follow failed", err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"target_type": body.TargetType, "target_id": body.TargetID, "following": r.Method != http.MethodDelete})
}

func (h *Handler) ReportPost(w http.ResponseWriter, r *http.Request) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	postID, err := uuid.Parse(r.PathValue("postId"))
	if err != nil {
		writeCommunityError(w, r, 400, "Invalid post ID", "postId must be a UUID")
		return
	}
	var body struct {
		Reason  string  `json:"reason"`
		Details *string `json:"details"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		writeCommunityError(w, r, 422, "Invalid report", "reason is required")
		return
	}
	var id uuid.UUID
	err = h.service.repo.pool.QueryRow(r.Context(), `INSERT INTO post_reports(post_id,principal_id,reason,details) VALUES($1,$2,$3,$4) ON CONFLICT(post_id,principal_id,reason) DO UPDATE SET details=EXCLUDED.details RETURNING report_id`, postID, viewer, body.Reason, body.Details).Scan(&id)
	if err != nil {
		writeCommunityError(w, r, 400, "Report failed", err.Error())
		return
	}
	writeJSON(w, 201, map[string]any{"report_id": id, "status": "pending"})
}

func (h *Handler) BlockPrincipal(w http.ResponseWriter, r *http.Request) {
	viewer, err := viewerUUID(r)
	if err != nil {
		writeCommunityError(w, r, 401, "Unauthorized", err.Error())
		return
	}
	blocked, err := uuid.Parse(r.PathValue("principalId"))
	if err != nil {
		writeCommunityError(w, r, 400, "Invalid principal ID", "principalId must be a UUID")
		return
	}
	if r.Method == http.MethodDelete {
		_, err = h.service.repo.pool.Exec(r.Context(), `DELETE FROM principal_blocks WHERE blocker_id=$1 AND blocked_id=$2`, viewer, blocked)
	} else {
		_, err = h.service.repo.pool.Exec(r.Context(), `INSERT INTO principal_blocks(blocker_id,blocked_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, viewer, blocked)
	}
	if err != nil {
		writeCommunityError(w, r, 400, "Block failed", err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"principal_id": blocked, "blocked": r.Method != http.MethodDelete})
}

func parseLimit(r *http.Request, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || value < 1 {
		return fallback
	}
	if value > 100 {
		return 100
	}
	return value
}
func viewerUUID(r *http.Request) (uuid.UUID, error) {
	id := httpserver.GetUserID(r.Context())
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("verified account or demo identity required")
	}
	return parsed, nil
}
func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeCommunityError(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
	httpserver.WriteError(w, status, title, detail, r.URL.Path, httpserver.GetRequestID(r.Context()))
}

var _ = pgx.ErrNoRows
