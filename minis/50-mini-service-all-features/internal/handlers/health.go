package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/database"
)

// Health returns a simple health check
func Health(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}
}

// Ready checks if the service is ready (database connected, etc.)
func Ready(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database
		if err := db.Health(r.Context()); err != nil {
			logger.Error().Err(err).Msg("database health check failed")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "not ready",
				"reason": "database unavailable",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ready",
		})
	}
}
