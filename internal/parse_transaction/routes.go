package parsetransaction

import (
	"expent-backend/configs"
	"expent-backend/internal/parse_transaction/handler"
	"expent-backend/internal/parse_transaction/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	svc := service.NewService(configs.AppConfig.GEMINI_API_KEY)
	h := handler.NewHandler(svc)

	rg.POST("/parse-transaction", h.ParseTransaction)
}
