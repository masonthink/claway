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

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Global middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Handlers
	ideaH := handler.NewIdeaHandler(svc)
	taskH := handler.NewTaskHandler(svc)
	docH := handler.NewDocumentHandler(svc)
	creditH := handler.NewCreditHandler(svc)
	proxyH := handler.NewProxyHandler(svc)
	authH := handler.NewAuthHandler(svc)
	computeH := handler.NewComputeHandler(svc)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Public routes (no auth required)
	v1.GET("/auth/x", authH.XLogin)
	v1.GET("/auth/x/callback", authH.XCallback)
	v1.GET("/auth/openclaw/callback", authH.OpenClawCallback) // legacy
	v1.GET("/ideas", ideaH.ListIdeas)
	v1.GET("/ideas/:id", ideaH.GetIdea)
	v1.GET("/ideas/:id/tasks", taskH.ListTasks)
	v1.GET("/tasks/:id", taskH.GetTask)
	v1.GET("/ideas/:id/compute", computeH.GetIdeaCompute)
	v1.GET("/platform/compute", computeH.GetPlatformCompute)

	// Auth-protected routes
	auth := v1.Group("", middleware.RequireAuth(cfg.JWTSecret))

	auth.GET("/auth/me", authH.GetMe)

	// Ideas (write operations)
	auth.POST("/ideas", ideaH.CreateIdea)
	auth.GET("/ideas/:id/context", ideaH.GetIdeaContext)
	auth.POST("/tasks/:id/claim", taskH.ClaimTask)
	auth.DELETE("/tasks/:id/claim", taskH.UnclaimTask)
	auth.POST("/tasks/:id/submit", taskH.SubmitTask)
	auth.POST("/tasks/:id/review", taskH.ReviewTask)

	// Documents
	auth.GET("/tasks/:id/document", docH.GetDocument)
	auth.GET("/tasks/:id/document/versions", docH.ListVersions)
	auth.GET("/tasks/:id/document/versions/:ver", docH.GetVersion)
	auth.PUT("/tasks/:id/document", docH.UpdateDocument)

	// PRD
	auth.POST("/ideas/:id/publish", docH.PublishPRD)

	// LLM Proxy
	auth.POST("/proxy/chat", proxyH.Chat)

	// Compute
	auth.GET("/me/compute", computeH.GetMyCompute)
	auth.GET("/me/compute/ideas/:id", computeH.GetMyIdeaCompute)
	auth.GET("/tasks/:id/compute", computeH.GetTaskCompute)

	// Credits
	auth.GET("/me/credits", creditH.GetMyCredits)
	auth.GET("/me/contributions", creditH.GetMyContributions)
	auth.POST("/prd/:id/purchase", creditH.PurchasePRD)

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
