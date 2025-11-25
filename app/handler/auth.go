package handler

import (
	"net/http"
	"strings"

	"github.com/brunobotter/chat-websocket/auth"
	"github.com/brunobotter/chat-websocket/dto"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {

	var cred dto.Auth

	if err := c.Bind(&cred); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	user := cred.User
	pass := cred.Password

	if user == "" || pass == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing user or password"})
	}

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
