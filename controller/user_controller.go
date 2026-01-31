package controller

import (
	"log/slog"
	"net/http"

	"golang-backend/dto"
	// "golang-backend/middleware"
	"golang-backend/service"
	"golang-backend/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// Login godoc
// @Summary      Login User
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.UserLoginRequest true "Login Credentials"
// @Success      200  {object} utils.Response{data=object{token=string}}
// @Failure      400  {object} utils.Response
// @Router       /login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var input dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.service.Login(input)
	if err != nil {
		utils.ErrorResponse(ctx, "Login Failed", http.StatusUnauthorized, err.Error())
		return
	}

	slog.Info("User logged in", "email", input.Email)
	utils.SuccessResponse(ctx, "Login Successful", gin.H{"token": token})
}

// Register godoc
// @Summary      Register User
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.UserRegisterRequest true "Register Data"
// @Success      201  {object} utils.Response{data=dto.UserResponse}
// @Failure      400  {object} utils.Response
// @Router       /register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var input dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := c.service.Register(input)
	if err != nil {
		utils.ErrorResponse(ctx, "Registration Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.CreatedResponse(ctx, "User Registered Successfully", userResponse)
}

// VerifyEmail godoc
// @Summary      Verify Email
// @Description  Verify user email with code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.VerifyEmailRequest true "Verification Code"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /verify-email [post]
func (c *UserController) VerifyEmail(ctx *gin.Context) {
	var input dto.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.VerifyEmail(input); err != nil {
		utils.ErrorResponse(ctx, "Verification Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Email Verified Successfully", nil)
}

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Send reset password code to email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.ForgotPasswordRequest true "Email"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /forgot-password [post]
func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var input dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ForgotPassword(input.Email); err != nil {
		utils.ErrorResponse(ctx, "Request Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Reset Code Sent Successfully", nil)
}

// ResetPassword godoc
// @Summary      Reset Password
// @Description  Reset password with code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.ResetPasswordRequest true "Reset Data"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ResetPassword(input); err != nil {
		utils.ErrorResponse(ctx, "Reset Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Password Reset Successfully", nil)
}

// ResendVerificationCode godoc
// @Summary      Resend Verification Code
// @Description  Resend email verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.ForgotPasswordRequest true "Email"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /resend-verification [post]
func (c *UserController) ResendVerificationCode(ctx *gin.Context) {
	var input dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ResendVerificationCode(input.Email); err != nil {
		utils.ErrorResponse(ctx, "Request Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Verification Code Sent Successfully", nil)
}

// ResendResetPasswordCode godoc
// @Summary      Resend Reset Password Code
// @Description  Resend reset password code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.ForgotPasswordRequest true "Email"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /resend-reset-code [post]
func (c *UserController) ResendResetPasswordCode(ctx *gin.Context) {
	var input dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ResendResetPasswordCode(input.Email); err != nil {
		utils.ErrorResponse(ctx, "Request Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Reset Code Sent Successfully", nil)
}

// Me godoc
// @Summary      Get Current User
// @Description  Get details of the currently logged-in user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} utils.Response
// @Failure      401  {object} utils.Response
// @Router       /me [get]
func (c *UserController) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, "Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	userResponse, err := c.service.GetMe(userID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, "Failed to fetch user", http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "User details", userResponse)
}

// GetUsers godoc
// @Summary      Get Users
// @Description  Get paginated list of users with search and filter
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(15)
// @Param        search query string false "Search term"
// @Param        is_verified query boolean false "Filter by verification status"
// @Success      200  {object} utils.PaginationResult
// @Failure      403  {object} utils.Response
// @Router       /users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	// Policy Check
	// Assuming module name is namespaced or simple "user"
	// if err := middleware.AuthorizeRead(ctx, "user"); err != nil {
	// 	utils.ErrorResponse(ctx, err.Error(), http.StatusForbidden, nil)
	// 	return
	// }

	search := ctx.Query("search")
	isVerified := ctx.Query("is_verified")
	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")

	filters := map[string]interface{}{
		"search":      search,
		"is_verified": isVerified,
		"sort_by":     sortBy,
		"sort_order":  sortOrder,
	}

	page, perPage := utils.GetPaginationParams(ctx)

	result, err := c.service.GetUsers(filters, page, perPage)
	if err != nil {
		utils.ErrorResponse(ctx, "Failed to fetch users", http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate execution time if needed, or just pass result.Pagination
	// The PHP example calculated it in the trait.
	// utils.PaginatedResponse expects (ctx, message, data, meta)
	// We can build meta here.

	meta := utils.BuildMeta(result.Pagination, 0) // Execution time 0 or calculate it
	utils.PaginatedResponse(ctx, "Users retrieved successfully", result.Items, meta)
}

// Setup2FA godoc
// @Summary      Setup 2FA
// @Description  Generate 2FA secret and QR code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} utils.Response{data=dto.Setup2FAResponse}
// @Failure      400  {object} utils.Response
// @Router       /2fa/setup [post]
func (c *UserController) Setup2FA(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, "Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	response, err := c.service.Setup2FA(userID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, "Failed to setup 2FA", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "2FA Setup Initiated", response)
}

// Verify2FA godoc
// @Summary      Verify 2FA
// @Description  Verify 2FA code and enable 2FA
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body dto.Verify2FARequest true "Verification Code"
// @Success      200  {object} utils.Response
// @Failure      400  {object} utils.Response
// @Router       /2fa/verify [post]
func (c *UserController) Verify2FA(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, "Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var input dto.Verify2FARequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, "Validation Failed", http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.Verify2FA(userID.(string), input.Code); err != nil {
		utils.ErrorResponse(ctx, "Verification Failed", http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "2FA Verified and Enabled", nil)
}
