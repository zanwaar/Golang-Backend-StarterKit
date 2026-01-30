package middleware

import (
	"net/http"
	"strings"

	"golang-backend/config"
	"golang-backend/entity"
	"golang-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		// Fetch user with roles and permissions
		var user entity.User
		if err := config.DB.Preload("Roles.Permissions").First(&user, "id = ?", userID).Error; err != nil {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}

		c.Set("currentUser", &user)
		c.Next()
	}
}
