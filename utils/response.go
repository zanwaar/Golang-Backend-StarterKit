package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func APIResponse(ctx *gin.Context, message string, statusCode int, data, errors, meta interface{}) {
	jsonResponse := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
		Errors:  errors,
		Meta:    meta,
	}

	ctx.JSON(statusCode, jsonResponse)
}

func SuccessResponse(ctx *gin.Context, message string, data interface{}) {
	// Meta handled by BuildMeta if needed, but for simple success it's nil
	APIResponse(ctx, message, http.StatusOK, data, nil, nil)
}

func CreatedResponse(ctx *gin.Context, message string, data interface{}) {
	APIResponse(ctx, message, http.StatusCreated, data, nil, nil)
}

func ErrorResponse(ctx *gin.Context, message string, statusCode int, errors interface{}) {
	APIResponse(ctx, message, statusCode, nil, errors, nil) // Meta usually nil for errors
}

func PaginatedResponse(ctx *gin.Context, message string, data, meta interface{}) {
	APIResponse(ctx, message, http.StatusOK, data, nil, meta)
}
