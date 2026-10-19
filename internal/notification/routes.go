package notification

import (
	deviceTokenRepository "expent-backend/internal/device_token/repository"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/middleware"
	"expent-backend/internal/notification/handler"
	"expent-backend/internal/notification/repository"
	"expent-backend/internal/notification/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers admin-only notification routes on the given
// router group. r should NOT already require user auth — AdminAuth is
// applied here so a regular user's JWT cannot be used to send notifications.
func RegisterRoutes(r *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	deviceTokenRepo := deviceTokenRepository.NewRepository(prismaClient)
	svc := service.NewService(repo, deviceTokenRepo)
	h := handler.NewHandler(svc)

	notificationGroup := r.Group("/notifications")
	notificationGroup.Use(middleware.AdminAuth())
	notificationGroup.POST("", h.Create)
}
