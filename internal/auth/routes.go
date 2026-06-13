package auth

import (
	"expent-backend/internal/auth/handler"
	"expent-backend/internal/auth/repository"
	"expent-backend/internal/auth/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up auth related endpoints.
func RegisterRoutes(r *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	authGroup := r.Group("/auth")
	authGroup.POST("/google", h.GoogleLogin)
	authGroup.POST("/test-login", h.TestLogin)
	authGroup.POST("/refresh", h.RefreshToken)
}
