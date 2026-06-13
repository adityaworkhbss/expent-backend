package transaction

import (
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/transaction/handler"
	"expent-backend/internal/transaction/repository"
	"expent-backend/internal/transaction/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	txGroup := rg.Group("/transactions")
	txGroup.GET("", h.ListTransactions)
	txGroup.POST("", h.CreateTransaction)
	txGroup.PUT("/:id", h.UpdateTransaction)
	txGroup.DELETE("/:id", h.DeleteTransaction)
}
