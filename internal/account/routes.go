package account

import (
    "github.com/gin-gonic/gin"
    "expent-backend/internal/account/handler"
    "expent-backend/internal/account/service"
    "expent-backend/internal/account/repository"
    "expent-backend/internal/infrastructure/prisma"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
    repo := repository.NewRepository(prismaClient)
    svc := service.NewService(repo)
    h := handler.NewHandler(svc)

    accGroup := rg.Group("/accounts")
    accGroup.GET("", h.ListAccounts)
    accGroup.POST("", h.CreateAccount)
    accGroup.DELETE("/:id", h.DeleteAccount)
}
