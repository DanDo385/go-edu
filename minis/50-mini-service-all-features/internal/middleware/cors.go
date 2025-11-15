package middleware

import (
	"net/http"
	"strings"

	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/config"
)

// CORS adds Cross-Origin Resource Sharing headers
func CORS(cfg config.CORSConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			if len(cfg.AllowedOrigins) > 0 {
				origin := cfg.AllowedOrigins[0]
				if origin == "*" || contains(cfg.AllowedOrigins, r.Header.Get("Origin")) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			if len(cfg.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods",
					strings.Join(cfg.AllowedMethods, ", "))
			}

			if len(cfg.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers",
					strings.Join(cfg.AllowedHeaders, ", "))
			}

			// Handle preflight request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
