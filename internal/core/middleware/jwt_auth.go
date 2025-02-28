package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"project-api/internal/core/common/utils"
	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func respondError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// isExcludedRoute checks if a path should bypass JWT authentication
func isExcludedRoute(path string) bool {
	excludedRoutes := []*regexp.Regexp{
		regexp.MustCompile(`^/$`),
		regexp.MustCompile(`^/api/v1/auth/.*$`),
	}
	for _, pattern := range excludedRoutes {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

func JWTAuthMiddleware(c *fiber.Ctx) error { // เปลี่ยน signature เป็น Fiber's Middleware Handler
	if isExcludedRoute(c.Path()) {
		return c.Next() // ข้าม middleware ถ้าเป็น excluded route
	}
	// Get the authorization header
	authHeader := c.Get("Authorization") // ใช้ c.Get() แทน r.Header.Get()
	if authHeader == "" {
		logger.Warn("Missing authorization header", zap.String("path", c.Path()))       // Log path context
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header") // ใช้ fiber.NewError เพื่อ return error
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		logger.Warn("Invalid authorization header format", zap.String("authHeader", authHeader), zap.String("path", c.Path())) // Log authHeader context
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")                                 // ใช้ fiber.NewError เพื่อ return error
	}
	tokenString := parts[1]
	token, err := jwt.ParseWithClaims(tokenString, &utils.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT.Signed), nil
	})
	if err != nil || !token.Valid {
		logger.Warn("invalid or expired token", zap.Error(err), zap.String("token", tokenString), zap.String("path", c.Path())) // Log token context
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")                                             // ใช้ fiber.NewError เพื่อ return error
	}
	claims, ok := token.Claims.(*utils.UserClaims)
	if !ok {
		logger.Error("invalid token claims", zap.String("token", tokenString), zap.String("path", c.Path())) // Log token and path context
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")                              // ใช้ fiber.NewError เพื่อ return error
	}
	// ใช้ c.Context() เพื่อเข้าถึง Go Context ของ Fiber
	ctx := context.WithValue(c.UserContext(), utils.GetUserContextKey(), claims)
	c.SetUserContext(ctx) // Set Go Context ลง Fiber Context

	return c.Next()
}
