package controller

import (
	"net/http"

	"golang-backend/dto"
	"golang-backend/entity"
	"golang-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{DB: db}
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role (Admin only)
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        role  body      dto.CreateRoleRequest  true  "Role Data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Router       /admin/roles [post]
func (rc *RoleController) CreateRole(c *gin.Context) {
	var input dto.CreateRoleRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest, nil)
		return
	}

	role := entity.Role{Name: input.Name}
	if err := rc.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, "Role already exists", http.StatusConflict, nil)
		return
	}

	utils.SuccessResponse(c, "Role created successfully", role)
}

// CreatePermission godoc
// @Summary      Create a new permission
// @Description  Create a new permission (Admin only)
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        permission  body      dto.CreatePermissionRequest  true  "Permission Data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Router       /admin/permissions [post]
func (rc *RoleController) CreatePermission(c *gin.Context) {
	var input dto.CreatePermissionRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest, nil)
		return
	}

	permission := entity.Permission{Name: input.Name}
	if err := rc.DB.Create(&permission).Error; err != nil {
		utils.ErrorResponse(c, "Permission already exists", http.StatusConflict, nil)
		return
	}

	utils.SuccessResponse(c, "Permission created successfully", permission)
}

// AssignRoleToUser godoc
// @Summary      Assign a role to a user
// @Description  Assign a role to a user (Admin only)
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        assignment  body      dto.AssignRoleRequest  true  "Assignment Data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /admin/assign-role [post]
func (rc *RoleController) AssignRoleToUser(c *gin.Context) {
	var input dto.AssignRoleRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest, nil)
		return
	}

	var user entity.User
	if err := rc.DB.Preload("Roles").First(&user, "id = ?", input.UserID).Error; err != nil {
		utils.ErrorResponse(c, "User not found", http.StatusNotFound, nil)
		return
	}

	var role entity.Role
	if err := rc.DB.Where("name = ?", input.Role).First(&role).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound, nil)
		return
	}

	// Check if already has role
	if user.HasRole(input.Role) {
		utils.ErrorResponse(c, "User already has this role", http.StatusBadRequest, nil)
		return
	}

	if err := rc.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		utils.ErrorResponse(c, "Failed to assign role", http.StatusInternalServerError, nil)
		return
	}

	utils.SuccessResponse(c, "Role assigned successfully", nil)
}

// AssignPermissionToRole godoc
// @Summary      Assign a permission to a role
// @Description  Assign a permission to a role (Admin only)
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        assignment  body      dto.AssignPermissionRequest  true  "Assignment Data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /admin/assign-permission [post]
func (rc *RoleController) AssignPermissionToRole(c *gin.Context) {
	var input dto.AssignPermissionRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest, nil)
		return
	}

	var role entity.Role
	if err := rc.DB.Preload("Permissions").Where("name = ?", input.RoleName).First(&role).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound, nil)
		return
	}

	var permission entity.Permission
	if err := rc.DB.Where("name = ?", input.PermissionName).First(&permission).Error; err != nil {
		utils.ErrorResponse(c, "Permission not found", http.StatusNotFound, nil)
		return
	}

	if err := rc.DB.Model(&role).Association("Permissions").Append(&permission); err != nil {
		utils.ErrorResponse(c, "Failed to assign permission", http.StatusInternalServerError, nil)
		return
	}

	utils.SuccessResponse(c, "Permission assigned to role successfully", nil)
}
