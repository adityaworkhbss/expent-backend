package emi

import (
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/emi/handler"
	"expent-backend/internal/emi/repository"
	"expent-backend/internal/emi/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	emiGroup := rg.Group("/emis")
	emiGroup.GET("", h.ListEmis)
	emiGroup.POST("", h.CreateEmi)
	emiGroup.PUT("/:id", h.UpdateEmi)
	emiGroup.DELETE("/:id", h.DeleteEmi)
}
