package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
    "go.uber.org/zap"
)

// Recovery middleware recovers from any panics and returns a standard error response.
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // Log the panic with stack trace
                zap.S().Error("panic recovered", zap.Any("error", r))
                // Return a generic error response
                shared.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
                c.Abort()
            }
        }()
        c.Next()
    }
}
