package handler

import (
	"net/http"
	"strings"

	"github.com/brunobotter/chat-websocket/internal/auth"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	user := c.FormValue("user")
	pass := c.FormValue("password")

	if pass != "1234" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	rooms := []string{"default", "vip"}
	access, _ := auth.GenerateAccessToken(user, rooms)
	refresh, _ := auth.GenerateRefreshToken(user)

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func Refresh(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := auth.ValidateRefreshToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid refresh token"})
	}

	rooms := []string{"default", "vip"}
	newAccess, _ := auth.GenerateAccessToken(user, rooms)

	return c.JSON(http.StatusOK, echo.Map{"access_token": newAccess})
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := auth.ValidateAccessToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		return next(c)
	}
}
