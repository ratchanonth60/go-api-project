package utils

import "context"

type contextKey string

const (
	userContextKey contextKey = "user"
)

func GetUserContextKey() contextKey {
	return userContextKey
}
func GetUserIDFromContext(ctx context.Context) (*UserClaims, bool) {
	valFromContext := ctx.Value(GetUserContextKey())
	userClaims, ok := valFromContext.(*UserClaims)
	return userClaims, ok
}
