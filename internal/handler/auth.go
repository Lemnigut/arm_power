package handler

import (
	"errors"
	"net/http"

	"arm_back/internal/model"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	tokens, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrConflict) {
			c.JSON(http.StatusConflict, model.ErrorResponse{Error: "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "login failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	tokens, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
