package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/config"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/database"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/handlers"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/metrics"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Setup logger
	logger := setupLogger(&cfg.Logging)

	logger.Info().Msg("Starting microservice...")

	// Setup metrics
	m := metrics.New()

	// Setup database
	db := database.New()
	defer db.Close()

	logger.Info().Msg("Database initialized")

	// Setup router
	router := setupRouter(cfg, db, logger, m)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info().Msgf("Server starting on %s", cfg.Server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("server shutdown failed")
	}

	logger.Info().Msg("Server stopped gracefully")
}

func setupLogger(cfg *config.LoggingConfig) zerolog.Logger {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	var logger zerolog.Logger
	if cfg.Format == "console" {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return logger
}

func setupRouter(
	cfg *config.Config,
	db *database.DB,
	logger zerolog.Logger,
	m *metrics.Metrics,
) http.Handler {
	mux := http.NewServeMux()

	// Health endpoints (no auth required)
	mux.HandleFunc("/health", handlers.Health(logger))
	mux.HandleFunc("/ready", handlers.Ready(db, logger))
	mux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	// Auth endpoints (no auth required)
	mux.HandleFunc("/login", handlers.Login(db, cfg.JWT, logger))

	// Protected endpoints (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/users", handlers.ListUsers(db, logger))
	protectedMux.HandleFunc("/users/", handlers.GetUser(db, logger))

	// Apply auth middleware to protected routes
	protected := middleware.Chain(
		protectedMux,
		middleware.Auth(cfg.JWT.Secret),
	)

	mux.Handle("/users", protected)
	mux.Handle("/users/", protected)

	// Apply global middleware to all routes
	handler := middleware.Chain(
		mux,
		middleware.Recovery(logger),
		middleware.RequestID(),
		middleware.Logging(logger),
		middleware.Metrics(m),
		middleware.CORS(cfg.CORS),
		middleware.RateLimit(cfg.RateLimit),
	)

	return handler
}
