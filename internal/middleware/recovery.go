package middleware

import (
	"expent-backend/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {

				zap.S().Error("panic recovered", zap.Any("error", r))

				shared.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
