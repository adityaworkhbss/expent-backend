package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "expent-backend/configs"
    "go.uber.org/zap"
)

// Auth middleware validates JWT and sets userId in context
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing Authorization header"})
            return
        }
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid Authorization header"})
            return
        }
        tokenString := parts[1]
        // Parse token
        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrTokenUnverifiable
            }
            return []byte(configs.AppConfig.JWT_SECRET), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
            return
        }
        // Extract userId claim (assume "sub" claim)
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token claims"})
            return
        }
        sub, ok := claims["sub"].(string)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token missing sub claim"})
            return
        }
        // Store userId in context for downstream handlers
        c.Set("userId", sub)
        c.Next()
    }
}
