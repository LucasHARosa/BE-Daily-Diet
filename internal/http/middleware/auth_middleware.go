package middleware

import (
	"context"
	"net/http"
	"strings"

	infraauth "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(jwtService *infraauth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				responses.Error(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			claims, err := jwtService.Validate(parts[1])
			if err != nil {
				responses.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}
