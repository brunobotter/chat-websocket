package handler

import (
	"net/http"
	"strings"

	"github.com/brunobotter/chat-websocket/internal/auth"
	"github.com/brunobotter/chat-websocket/internal/dto"
	"github.com/labstack/echo/v4"
)

// AuthService defines the interface for authentication operations
// This allows for easier mocking and testing
//go:generate mockgen -destination=../../mocks/mock_auth_service.go -package=mocks github.com/brunobotter/chat-websocket/internal/handler AuthService

type AuthService interface {
	GenerateAccessToken(user string, rooms []string) (string, error)
	GenerateRefreshToken(user string) (string, error)
	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, error)
}

type AuthHandler struct {
	AuthSvc AuthService
}

func NewAuthHandler(authSvc AuthService) *AuthHandler {
	return &AuthHandler{AuthSvc: authSvc}
}

func (h *AuthHandler) Login(c echo.Context) error {
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
	access, err := h.AuthSvc.GenerateAccessToken(user, rooms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate access token"})
	}
	refresh, err := h.AuthSvc.GenerateRefreshToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate refresh token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := h.AuthSvc.ValidateRefreshToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid refresh token"})
	}

	rooms := []string{"default", "vip"}
	newAccess, err := h.AuthSvc.GenerateAccessToken(user, rooms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate access token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": newAccess})
}

func (h *AuthHandler) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := h.AuthSvc.ValidateAccessToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		return next(c)
	}
}
