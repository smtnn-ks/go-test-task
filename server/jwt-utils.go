package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTokens(userId string) (accessToken string, refreshToken string, err error) {

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().UTC().Add(time.Minute * 15).Unix(),
	}).SignedString([]byte(secrets.jwtAccessSecret))

	if err != nil {
		return
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().UTC().Add(time.Hour * 24 * 7).Unix(),
	}).SignedString([]byte(secrets.jwtRefreshSecret))

	if err != nil {
		return
	}

	return
}
