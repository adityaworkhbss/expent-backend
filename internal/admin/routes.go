package admin

import (
	"expent-backend/internal/admin/handler"
	"expent-backend/internal/admin/repository"
	"expent-backend/internal/admin/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the public admin routes (login and the one-time
// setup endpoint). Both are public because there is no admin session yet at
// the point they're called.
func RegisterRoutes(r *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	adminGroup := r.Group("/admin")
	adminGroup.POST("/login", h.Login)
	adminGroup.POST("/setup", h.Setup)
}
