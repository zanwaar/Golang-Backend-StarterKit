package utils

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaginationResult struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	Total        int64 `json:"total"`
	LastPage     int   `json:"last_page"`
	From         int   `json:"from"`
	To           int   `json:"to"`
	HasMorePages bool  `json:"has_more_pages"`
}

func BuildMeta(pagination Pagination, executionTime float64) map[string]interface{} {
	return map[string]interface{}{
		"pagination":     pagination,
		"execution_time": executionTime,
		"timestamp":      time.Now().Format(time.RFC3339),
	}
}

func GetPaginationParams(ctx *gin.Context) (page, perPage int) {
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("per_page", "15")

	page, _ = strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPage, _ = strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}

	return page, perPage
}

func CalculatePagination(total int64, page, perPage int) Pagination {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	from := (page-1)*perPage + 1
	to := from + perPage - 1
	if int64(to) > total {
		to = int(total)
	}
	if int64(from) > total { // Case where page is out of range
		from = 0 // Or handle differently
		to = 0
	}

	return Pagination{
		CurrentPage:  page,
		PerPage:      perPage,
		Total:        total,
		LastPage:     lastPage,
		From:         from,
		To:           to,
		HasMorePages: page < lastPage,
	}
}
