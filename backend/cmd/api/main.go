package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/content"
	"github.com/tourtect/backend/internal/places"
	"github.com/tourtect/backend/internal/platform/config"
	"github.com/tourtect/backend/internal/platform/database"
	"github.com/tourtect/backend/internal/platform/httpserver"
	"github.com/tourtect/backend/internal/platform/logging"
	"github.com/tourtect/backend/internal/pricing"
	"github.com/tourtect/backend/internal/safety"
)

type Server struct {
	db            *database.DB
	placesHandler *places.Handler
	postsHandler  *content.Handler
	priceHandler  *pricing.Handler
	safetyHandler *safety.Handler
}

// Ensure Server implements openapi.ServerInterface
var _ openapi.ServerInterface = (*Server)(nil)

func (s *Server) HealthLive(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	resp := openapi.HealthStatus{
		Status:    openapi.Ok,
		Timestamp: &now,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) HealthReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := make(map[string]string)
	status := openapi.Ok

	if err := s.db.Ping(ctx); err != nil {
		status = openapi.Unavailable
		checks["postgres"] = "DOWN: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		checks["postgres"] = "UP"
	}

	now := time.Now()
	resp := openapi.HealthStatus{
		Status:    status,
		Checks:    &checks,
		Timestamp: &now,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) ListPlaces(w http.ResponseWriter, r *http.Request, params openapi.ListPlacesParams) {
	s.placesHandler.List(w, r, params)
}

func (s *Server) GetPlace(w http.ResponseWriter, r *http.Request, placeId openapi.PlaceIdParam, params openapi.GetPlaceParams) {
	s.placesHandler.Get(w, r, placeId, params)
}

func (s *Server) ListPosts(w http.ResponseWriter, r *http.Request, params openapi.ListPostsParams) {
	s.postsHandler.List(w, r, params)
}

func (s *Server) CreateDraft(w http.ResponseWriter, r *http.Request, params openapi.CreateDraftParams) {
	s.postsHandler.CreateDraft(w, r, params)
}

func (s *Server) PublishPost(w http.ResponseWriter, r *http.Request, postId openapi.PostIdParam, params openapi.PublishPostParams) {
	s.postsHandler.Publish(w, r, postId, params)
}

func (s *Server) GetFeed(w http.ResponseWriter, r *http.Request, _ openapi.GetFeedParams) {
	s.postsHandler.Feed(w, r)
}
func (s *Server) SearchCommunity(w http.ResponseWriter, r *http.Request, _ openapi.SearchCommunityParams) {
	s.postsHandler.Search(w, r)
}
func (s *Server) ListComments(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.ListCommentsParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.Comments(w, r)
}
func (s *Server) CreateComment(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.CreateCommentParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.Comments(w, r)
}
func (s *Server) MarkPostUseful(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.MarkPostUsefulParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.UsefulVote(w, r)
}
func (s *Server) UnmarkPostUseful(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.UnmarkPostUsefulParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.UsefulVote(w, r)
}
func (s *Server) ListSavedPosts(w http.ResponseWriter, r *http.Request, _ openapi.ListSavedPostsParams) {
	s.postsHandler.SavedList(w, r)
}
func (s *Server) SavePost(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.SavePostParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.SavedPost(w, r)
}
func (s *Server) UnsavePost(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.UnsavePostParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.SavedPost(w, r)
}
func (s *Server) ListNotifications(w http.ResponseWriter, r *http.Request, _ openapi.ListNotificationsParams) {
	s.postsHandler.Notifications(w, r)
}
func (s *Server) UpdateNotifications(w http.ResponseWriter, r *http.Request, _ openapi.UpdateNotificationsParams) {
	s.postsHandler.Notifications(w, r)
}
func (s *Server) CreateFollow(w http.ResponseWriter, r *http.Request, _ openapi.CreateFollowParams) {
	s.postsHandler.Follow(w, r)
}
func (s *Server) DeleteFollow(w http.ResponseWriter, r *http.Request, _ openapi.DeleteFollowParams) {
	s.postsHandler.Follow(w, r)
}
func (s *Server) ReportPost(w http.ResponseWriter, r *http.Request, postID openapi.PostIdParam, _ openapi.ReportPostParams) {
	r.SetPathValue("postId", postID.String())
	s.postsHandler.ReportPost(w, r)
}
func (s *Server) BlockPrincipal(w http.ResponseWriter, r *http.Request, principalID openapi_types.UUID, _ openapi.BlockPrincipalParams) {
	r.SetPathValue("principalId", principalID.String())
	s.postsHandler.BlockPrincipal(w, r)
}
func (s *Server) UnblockPrincipal(w http.ResponseWriter, r *http.Request, principalID openapi_types.UUID, _ openapi.UnblockPrincipalParams) {
	r.SetPathValue("principalId", principalID.String())
	s.postsHandler.BlockPrincipal(w, r)
}

func (s *Server) CreatePriceCheck(w http.ResponseWriter, r *http.Request, params openapi.CreatePriceCheckParams) {
	s.priceHandler.CreatePriceCheck(w, r, params)
}

func (s *Server) CreateSafetyAssessment(w http.ResponseWriter, r *http.Request, params openapi.CreateSafetyAssessmentParams) {
	s.safetyHandler.CreateSafetyAssessment(w, r, params)
}

func main() {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	logging.Init(cfg.LogLevel)
	slog.Info("Starting Tourtect API Server", slog.String("port", cfg.Port))

	// 3. Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// 4. Instantiate repositories and services
	placesRepo := places.NewRepository(db.Pool)
	placesService := places.NewService(placesRepo)
	placesHandler := places.NewHandler(placesService)

	postsRepo := content.NewRepository(db.Pool)
	postsService := content.NewService(postsRepo)
	postsHandler := content.NewHandler(postsService)

	priceEngine := pricing.NewEngine(db.Pool)
	priceHandler := pricing.NewHandler(priceEngine)

	safetyEngine := safety.NewEngine(db.Pool)
	safetyHandler := safety.NewHandler(safetyEngine)

	// 5. Instantiate server
	apiServer := &Server{
		db:            db,
		placesHandler: placesHandler,
		postsHandler:  postsHandler,
		priceHandler:  priceHandler,
		safetyHandler: safetyHandler,
	}

	// 6. Generate oapi router handler
	openapiHandler := openapi.Handler(apiServer)

	// 7. Wire standard middlewares
	mux := http.NewServeMux()
	mux.Handle("/", openapiHandler)

	var handler http.Handler = mux
	handler = httpserver.AuthBoundary(handler)
	handler = httpserver.Logging(handler)
	handler = httpserver.RequestID(handler)
	handler = httpserver.PanicRecovery(handler)
	handler = httpserver.CORS(handler)
	handler = httpserver.SecurityHeaders(handler)
	handler = httpserver.BodySizeLimit(10 * 1024 * 1024)(handler) // 10MB limit
	handler = httpserver.Timeout(30 * time.Second)(handler)

	// 8. Start HTTP Server with graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 35 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown coordination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("Shutting down API server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed, forcing close", slog.Any("error", err))
	} else {
		slog.Info("Server stopped successfully")
	}
}
