package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/standard-group/mesa/internal/jwt"
)

type contextKey string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Warn().Msg("Missing token in request")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing token"})
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid token")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
			return
		}

		log.Debug().Str("user_id", claims.UserID).Msg("Authenticated request")
		ctx := context.WithValue(r.Context(), contextKey("user_id"), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
