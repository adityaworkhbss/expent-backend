package device_token

import (
	"expent-backend/internal/device_token/handler"
	"expent-backend/internal/device_token/repository"
	"expent-backend/internal/device_token/service"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	deviceTokenGroup := r.Group("/device-tokens")
	deviceTokenGroup.Use(middleware.Auth())
	deviceTokenGroup.POST("", h.RegisterDeviceToken)
}