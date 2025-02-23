package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

func respondError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Missing authorization header")
			respondError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Invalid authorization header format")
			respondError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}
		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config.JWT.Signed), nil
		})
		if err != nil || !token.Valid {
			logger.Error("invalid or expired token", zap.Error(err))
			respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			logger.Error("invalid token claims")
			respondError(w, http.StatusUnauthorized, "invalid token claims")
		}
		ctx := context.WithValue(r.Context(), userContextKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userContextKey).(string)
	return v, ok
}
