package dashboard

import (
	"expent-backend/internal/dashboard/handler"
	"expent-backend/internal/dashboard/repository"
	"expent-backend/internal/dashboard/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)
	dashGroup := rg.Group("/dashboard")
	dashGroup.GET("", h.GetDashboard)
}
