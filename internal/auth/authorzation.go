package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User  string   `json:"user"`
	Rooms []string `json:"rooms"`
	jwt.RegisteredClaims
}

var (
	accessSecret  = []byte("super-secret-access-key") // Idealmente viria do .env
	refreshSecret = []byte("super-secret-refresh-key")
	AccessTTL     = 5 * time.Minute
	RefreshTTL    = 24 * time.Hour
)

// GenerateAccessToken cria o JWT de acesso
func GenerateAccessToken(user string, rooms []string) (string, error) {
	claims := Claims{
		User:  user,
		Rooms: rooms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chat-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken cria o JWT de refresh (sem rooms)
func GenerateRefreshToken(user string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chat-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateAccessToken valida o token de acesso e retorna o usuário e as rooms
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateRefreshToken valida o token de refresh e retorna o usuário
func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", errors.New("invalid refresh token")
}
