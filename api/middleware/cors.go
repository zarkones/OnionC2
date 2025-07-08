package middleware

import (
	"api/config"
	"net/http"
)

func WithCORS(h http.Handler) http.Handler {

	if *config.AllowedOrigin == "" {
		// This disables CORS.
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", *config.AllowedOrigin)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
