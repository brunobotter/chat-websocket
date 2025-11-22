package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager defines the interface for token operations, improving testability and flexibility.
type TokenManager interface {
	GenerateAccessToken(user string, rooms []string) (string, error)
	GenerateRefreshToken(user string) (string, error)
	ValidateAccessToken(tokenStr string) (*Claims, error)
	ValidateRefreshToken(tokenStr string) (string, error)
}

type Claims struct {
	User  string   `json:"user"`
	Rooms []string `json:"rooms"`
	jwt.RegisteredClaims
}

var (
	accessSecret  = []byte("132")
	refreshSecret = []byte("123")
	AccessTTL     = 5 * time.Minute
	RefreshTTL    = 24 * time.Hour
)

type JWTManager struct{}

func NewJWTManager() *JWTManager {
	return &JWTManager{}
}

func (j *JWTManager) GenerateAccessToken(user string, rooms []string) (string, error) {
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

func (j *JWTManager) GenerateRefreshToken(user string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chat-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func (j *JWTManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
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

func (j *JWTManager) ValidateRefreshToken(tokenStr string) (string, error) {
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
