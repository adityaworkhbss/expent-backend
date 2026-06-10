package shared

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse returns a standard success JSON response.
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(200, APIResponse{Success: true, Message: message, Data: data})
}

// ErrorResponse returns a standard error JSON response with given HTTP status.
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{Success: false, Message: message})
}
