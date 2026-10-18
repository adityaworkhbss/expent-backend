package parsetransaction

import (
	"expent-backend/configs"
	accountrepo "expent-backend/internal/account/repository"
	categoryrepo "expent-backend/internal/category/repository"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/parse_transaction/handler"
	"expent-backend/internal/parse_transaction/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	svc := service.NewService(
		configs.AppConfig.GEMINI_API_KEY,
		configs.AppConfig.GEMINI_MODEL,
		categoryrepo.NewRepository(prismaClient),
		accountrepo.NewRepository(prismaClient),
	)
	h := handler.NewHandler(svc)

	rg.POST("/parse-transaction", h.ParseTransaction)
}
