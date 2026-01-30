package middleware

import (
	"errors"
	"golang-backend/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Policy helpers to mimic Laravel's Authorization policies
// In Go, we can't easily use Traits, but we can use helper functions.

func Authorize(ctx *gin.Context, action, module string) error {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		return errors.New("unauthorized")
	}

	user, ok := currentUser.(*entity.User)
	if !ok {
		return errors.New("invalid user context")
	}

	permissionName := action + "_" + module
	if !user.HasPermission(permissionName) {
		return errors.New("forbidden: insufficient permissions")
	}

	return nil
}

func AuthorizeRead(ctx *gin.Context, module string) error {
	// Usually maps to 'view_{module}' or 'read_{module}'
	// Let's assume 'view' based on common conventions or 'read'
	// The user example had authorizeRead. Let's use 'view' as a common standard or 'read'.
	// If the user's PHP system uses 'admin_unit' module, it likely has 'view_admin_unit'.
	// I'll stick to 'view' or 'read'. Let's use 'view' for UI friendliness or 'read' for CRUD.
	// Let's genericize to just pass the action. But to match PHP 'authorizeRead', I'll make a choice.
	// PHP policy usually maps method to permission.
	return Authorize(ctx, "view", module)
}

func AuthorizeCreate(ctx *gin.Context, module string) error {
	return Authorize(ctx, "create", module)
}

func AuthorizeEdit(ctx *gin.Context, module string) error {
	return Authorize(ctx, "edit", module)
}

func AuthorizeDelete(ctx *gin.Context, module string) error {
	return Authorize(ctx, "delete", module)
}

// Helper to handle the error response automatically if desired, similar to Abort API
func EnsurePermission(ctx *gin.Context, action, module string) bool {
	if err := Authorize(ctx, action, module); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": err.Error(),
		})
		ctx.Abort()
		return false
	}
	return true
}
