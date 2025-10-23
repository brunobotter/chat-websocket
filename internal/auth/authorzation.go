package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtAccessSecret  = []byte("sua_chave_secreta_access")
	jwtRefreshSecret = []byte("sua_chave_secreta_refresh")
)

type Claims struct {
	User  string   `json:"user"`
	Rooms []string `json:"rooms"`
	jwt.RegisteredClaims
}

// ------------------- ACCESS TOKEN -------------------

// Gera access token (curta duração)
func GenerateAccessToken(user string, rooms []string) (string, error) {
	claims := Claims{
		User:  user,
		Rooms: rooms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAccessSecret)
}

// Valida access token
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtAccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// ------------------- REFRESH TOKEN -------------------

// Gera refresh token (24h)
func GenerateRefreshToken(user string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtRefreshSecret)
}

// Valida refresh token
func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("refresh token inválido")
	}

	return claims.Subject, nil
}
