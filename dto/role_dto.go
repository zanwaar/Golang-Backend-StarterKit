package dto

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type AssignPermissionRequest struct {
	RoleName       string `json:"role_name" binding:"required"`
	PermissionName string `json:"permission_name" binding:"required"`
}
