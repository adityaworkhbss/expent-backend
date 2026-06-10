package budget

import (
	"expent-backend/internal/budget/handler"
	"expent-backend/internal/budget/repository"
	"expent-backend/internal/budget/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	budGroup := rg.Group("/budgets")
	budGroup.GET("", h.ListBudgets)
	budGroup.POST("", h.CreateBudget)
	budGroup.DELETE("/:id", h.DeleteBudget)
}
