package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// Logger middleware logs each request using zap.
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Stop timer
        latency := time.Since(start)
        if raw != "" {
            path = path + "?" + raw
        }
        // Use zap's global logger
        zap.S().Info(
            "Incoming request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("client_ip", c.ClientIP()),
        )
    }
}
