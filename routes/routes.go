package routes

import (
	"golang-backend/controller"
	"golang-backend/middleware"
	"golang-backend/utils"

	_ "golang-backend/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(app *gin.Engine, userCtrl *controller.UserController, roleCtrl *controller.RoleController) {
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := app.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			utils.SuccessResponse(c, "pong", nil)
		})
		api.POST("/register", userCtrl.Register)
		api.POST("/login", userCtrl.Login)
		api.POST("/verify-email", userCtrl.VerifyEmail)
		api.POST("/forgot-password", userCtrl.ForgotPassword)
		api.POST("/reset-password", userCtrl.ResetPassword)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", userCtrl.Me)
			protected.GET("/users", userCtrl.GetUsers)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.RoleAuthMiddleware("admin"))
		{
			admin.POST("/roles", roleCtrl.CreateRole)
			admin.POST("/permissions", roleCtrl.CreatePermission)
			admin.POST("/assign-role", roleCtrl.AssignRoleToUser)
			admin.POST("/assign-permission", roleCtrl.AssignPermissionToRole)
		}
	}
}
