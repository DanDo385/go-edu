package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/config"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/database"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/middleware"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/models"
)

// Login authenticates a user and returns a JWT token
func Login(db *database.DB, cfg config.JWTConfig, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse request
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"INVALID_REQUEST",
				"Invalid request body",
			))
			return
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"MISSING_CREDENTIALS",
				"Username and password are required",
			))
			return
		}

		// Authenticate
		user, err := db.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			logger.Warn().
				Str("username", req.Username).
				Msg("authentication failed")

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"INVALID_CREDENTIALS",
				"Invalid username or password",
			))
			return
		}

		// Generate JWT token
		token, err := middleware.GenerateJWT(
			user.ID,
			user.Username,
			cfg.Secret,
			cfg.Expiration,
		)
		if err != nil {
			logger.Error().Err(err).Msg("failed to generate JWT")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"INTERNAL_ERROR",
				"Failed to generate token",
			))
			return
		}

		// Return token
		expiresAt := time.Now().Add(cfg.Expiration).Unix()
		resp := models.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
		}

		logger.Info().
			Str("username", user.Username).
			Int("user_id", user.ID).
			Msg("user logged in")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
