package utils

import (
	"time"

	"project-api/internal/infra/config"

	"github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	AccessExp    *jwt.NumericDate `json:"access_exp"`
	RefreshExp   *jwt.NumericDate `json:"refresh_exp"`
}

func GenerateJWT(email string) (*TokenDetails, error) {

	signed := []byte(config.Config.JWT.Signed)
	td := &TokenDetails{
		AccessExp:  jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		RefreshExp: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   td.AccessExp,
	},
	)
	at, err := accessToken.SignedString(signed)
	if err != nil {
		return nil, err
	}
	td.AccessToken = at
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   td.RefreshExp,
	})
	rt, err := refreshToken.SignedString(signed)
	if err != nil {
		return nil, err
	}
	td.RefreshToken = rt
	return td, nil
}
