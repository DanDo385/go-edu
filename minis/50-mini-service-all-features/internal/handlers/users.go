package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/database"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/middleware"
	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/models"
)

// ListUsers returns all users
func ListUsers(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated user from context
		user, ok := middleware.GetUser(r.Context())
		if !ok {
			logger.Error().Msg("user not found in context")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"UNAUTHORIZED",
				"Authentication required",
			))
			return
		}

		logger.Debug().
			Str("username", user.Username).
			Int("user_id", user.UserID).
			Msg("listing users")

		// Query database
		users, err := db.ListUsers(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("failed to list users")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"INTERNAL_ERROR",
				"Failed to retrieve users",
			))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

// GetUser returns a single user by ID
func GetUser(db *database.DB, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated user from context
		authUser, ok := middleware.GetUser(r.Context())
		if !ok {
			logger.Error().Msg("user not found in context")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"UNAUTHORIZED",
				"Authentication required",
			))
			return
		}

		// Extract user ID from URL
		// URL format: /users/{id}
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"INVALID_USER_ID",
				"Invalid user ID",
			))
			return
		}

		logger.Debug().
			Str("username", authUser.Username).
			Int("user_id", authUser.UserID).
			Int("requested_user_id", id).
			Msg("getting user")

		// Query database
		user, err := db.GetUser(r.Context(), id)
		if err != nil {
			logger.Warn().
				Err(err).
				Int("user_id", id).
				Msg("user not found")

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"USER_NOT_FOUND",
				"User not found",
			))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
