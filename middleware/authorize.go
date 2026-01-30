package middleware

import (
	"net/http"

	"golang-backend/entity"
	"golang-backend/utils"

	"github.com/gin-gonic/gin"
)

func RoleAuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("currentUser")
		if !exists {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		user := currentUser.(*entity.User)

		// If user has 'admin' role, they can probably do anything, but let's be strict for now or allow it
		// For now, checks if user has ANY of the passed roles
		hasRole := false
		for _, role := range roles {
			if user.HasRole(role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.ErrorResponse(c, "Forbidden: Insufficient Role", http.StatusForbidden, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func PermissionAuthMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("currentUser")
		if !exists {
			utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		user := currentUser.(*entity.User)

		// Check if user has ANY of the required permissions
		hasPermission := false
		for _, perm := range permissions {
			if user.HasPermission(perm) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, "Forbidden: Insufficient Permissions", http.StatusForbidden, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
