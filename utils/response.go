package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseFormat -> konsisten dipakai untuk semua response API
type ResponseFormat struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondSuccess -> format untuk response sukses
func RespondSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, ResponseFormat{
		Status:  "success",
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// RespondError -> format untuk response error
func RespondError(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, ResponseFormat{
		Status:  "error",
		Code:    code,
		Message: message,
		Error:   err,
	})
}
