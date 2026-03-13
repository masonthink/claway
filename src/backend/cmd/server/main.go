package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/claway/server/internal/config"
	"github.com/claway/server/internal/handler"
	"github.com/claway/server/internal/middleware"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection pool
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	// Redis client (optional for MVP)
	if cfg.RedisURL != "" {
		redisOpts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Fatalf("failed to parse redis URL: %v", err)
		}
		rdb := redis.NewClient(redisOpts)
		defer rdb.Close()

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.Printf("warning: redis not available: %v (continuing without redis)", err)
		} else {
			log.Println("connected to redis")
		}
	} else {
		log.Println("redis not configured, skipping")
	}

	// Wire up Store -> Service -> Handlers
	st := store.New(dbpool)
	svc := service.New(st, cfg)

	// Start background cleanup for expired auth sessions
	bgCtx, bgCancel := context.WithCancel(context.Background())
	defer bgCancel()
	svc.StartAuthSessionCleanup(bgCtx)

	// Start reveal ticker (check for expired ideas every minute)
	svc.RunRevealTicker(bgCtx, 1*time.Minute)
	log.Println("reveal ticker started (1 min interval)")

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Global middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.RateLimiterWithConfig(echomw.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			// Skip rate limiting for health checks
			return c.Path() == "/health"
		},
		Store: echomw.NewRateLimiterMemoryStoreWithConfig(
			echomw.RateLimiterMemoryStoreConfig{Rate: 30, Burst: 60, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use X-Forwarded-For for Cloudflare/Caddy proxied requests
			if xff := c.Request().Header.Get("X-Forwarded-For"); xff != "" {
				return xff, nil
			}
			return c.RealIP(), nil
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded, please try again later",
			})
		},
	}))
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{
			"https://claway.cc",
			"https://www.claway.cc",
			"http://localhost:3000",
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Handlers
	ideaH := handler.NewIdeaHandler(svc)
	contribH := handler.NewContributionHandler(svc)
	voteH := handler.NewVoteHandler(svc)
	authH := handler.NewAuthHandler(svc)
	userH := handler.NewUserHandler(svc)
	statsH := handler.NewStatsHandler(svc)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// --- Public routes (no auth required) ---

	// Auth
	v1.GET("/auth/x", authH.XLogin)
	v1.GET("/auth/x/callback", authH.XCallback)
	v1.GET("/auth/openclaw/callback", authH.OpenClawCallback) // legacy
	v1.POST("/auth/session", authH.CreateAuthSession)         // agent session flow
	v1.GET("/auth/session/:sid", authH.GetAuthSession)        // agent session polling

	// Public API
	pub := v1.Group("/public")
	pub.GET("/stats", statsH.GetPlatformStats)
	pub.GET("/ideas", ideaH.ListIdeas)
	pub.GET("/ideas/:id", ideaH.GetIdea)
	pub.GET("/ideas/:id/contributions", contribH.ListContributions)
	pub.GET("/ideas/:id/result", ideaH.GetRevealResult)
	pub.GET("/users/:username", userH.GetUserProfile)

	// --- Auth-protected routes ---
	auth := v1.Group("", middleware.RequireAuth(cfg.JWTSecret))

	// User
	auth.GET("/auth/me", authH.GetMe)
	auth.GET("/me", authH.GetMe)

	// Ideas
	auth.POST("/ideas", ideaH.CreateIdea)
	auth.GET("/ideas", ideaH.ListIdeas)
	auth.GET("/ideas/:id", ideaH.GetIdea)
	auth.GET("/ideas/:id/result", ideaH.GetRevealResult)

	// Contributions
	auth.POST("/ideas/:id/contributions", contribH.CreateContribution)
	auth.GET("/ideas/:id/contributions", contribH.ListContributions)
	auth.PUT("/contributions/:id", contribH.UpdateContribution)
	auth.POST("/contributions/:id/submit", contribH.SubmitContribution)
	auth.GET("/contributions/:id", contribH.GetContribution)

	// Votes
	auth.POST("/ideas/:id/votes", voteH.CastVote)

	// My data
	auth.GET("/me/ideas", ideaH.ListMyIdeas)
	auth.GET("/me/contributions", contribH.ListMyContributions)
	auth.GET("/me/votes", voteH.ListMyVotes)

	// Draft preview (author only)
	auth.GET("/draft/:contribution_id", contribH.GetDraftPreview)

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}
