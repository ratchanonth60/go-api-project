package utils

import (
	"time"

	"project-api/internal/core/entity"
	"project-api/internal/infra/config"

	"github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	AccessExp    *jwt.NumericDate `json:"access_exp"`
	RefreshExp   *jwt.NumericDate `json:"refresh_exp"`
}
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"` // Expiration Time
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"` // Issued At Time
	Subject   string           `json:"sub,omitempty"` // Subject (User ID)
	Issuer    string           `json:"iss,omitempty"` // Issuer (Optional)
	Audience  []string         `json:"aud,omitempty"` // Audience (Optional)
	jwt.RegisteredClaims
}

func GenerateJWT(user *entity.User) (*TokenDetails, error) {

	signed := []byte(config.Config.JWT.Signed)
	td := &TokenDetails{
		AccessExp:  jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		RefreshExp: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		UserID:    user.ID,
		Username:  user.UserName,
		Email:     user.Email,
		ExpiresAt: td.AccessExp,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.UserName,
	})

	at, err := accessToken.SignedString(signed)
	if err != nil {
		return nil, err
	}
	td.AccessToken = at
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		UserID:    user.ID,
		Username:  user.UserName,
		Email:     user.Email,
		ExpiresAt: td.AccessExp,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.UserName,
	})
	rt, err := refreshToken.SignedString(signed)
	if err != nil {
		return nil, err
	}
	td.RefreshToken = rt
	return td, nil
}
