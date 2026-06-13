package usercustomization

import (
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/user_customization/handler"
	"expent-backend/internal/user_customization/repository"
	"expent-backend/internal/user_customization/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	custGroup := rg.Group("/user-customization")
	custGroup.GET("", h.GetCustomization)
	custGroup.PUT("", h.UpdateCustomization)
}
