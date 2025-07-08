package middleware

import (
	"api/config"
	"net/http"
	"strings"
)

func WithCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if config.AllowedOrigins != nil && *config.AllowedOrigins != "" {
			allowedList := strings.Split(*config.AllowedOrigins, ",")

			for _, allowed := range allowedList {
				allowed = strings.TrimSpace(allowed)
				if allowed == "*" {
					// Wildcard: allow any origin
					w.Header().Set("Access-Control-Allow-Origin", "*")
					break
				} else if allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
