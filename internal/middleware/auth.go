package middleware

import (
	"net/http"
	"strings"

	"arm_back/internal/model"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UserIDKey = "userID"

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid authorization format"})
			return
		}

		userID, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid token"})
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(UserIDKey)
	return id.(uuid.UUID)
}
