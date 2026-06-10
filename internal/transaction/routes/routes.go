package transaction

import (
    "github.com/gin-gonic/gin"
    "expent-backend/internal/transaction/handler"
    "expent-backend/internal/transaction/service"
    "expent-backend/internal/transaction/repository"
    "expent-backend/internal/infrastructure/prisma"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
    repo := repository.NewRepository(prismaClient)
    svc := service.NewService(repo)
    h := handler.NewHandler(svc)

    txGroup := rg.Group("/transactions")
    txGroup.GET("", h.ListTransactions)
    txGroup.POST("", h.CreateTransaction)
    txGroup.DELETE("/:id", h.DeleteTransaction)
}
