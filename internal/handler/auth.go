package handler

import (
	"net/http"
	"strings"

	"github.com/brunobotter/chat-websocket/internal/auth"
	"github.com/brunobotter/chat-websocket/internal/dto"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var cred dto.Auth
	if err := c.Bind(&cred); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if cred.User == "" || cred.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing user or password"})
	}

	if cred.Password != "1234" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	rooms := []string{"default", "vip"}
	access, err := auth.GenerateAccessToken(cred.User, rooms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate access token"})
	}
	refresh, err := auth.GenerateRefreshToken(cred.User)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate refresh token"})
	}

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
	newAccess, err := auth.GenerateAccessToken(user, rooms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate access token"})
	}

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
