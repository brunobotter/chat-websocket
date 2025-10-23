package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey          = []byte("supersecretkey123")
	jwtRefreshKey   = []byte("refreshsecretkey456")
	AccessTokenTTL  = 5 * time.Minute
	RefreshTokenTTL = 24 * time.Hour
)

type Claims struct {
	User  string   `json:"user"`
	Rooms []string `json:"rooms"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user string, rooms []string) (string, error) {
	claims := &Claims{
		User:  user,
		Rooms: rooms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(user string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   user,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtRefreshKey)
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtRefreshKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}
	return claims.Subject, nil
}
