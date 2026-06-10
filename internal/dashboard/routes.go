package dashboard

import (
	"expent-backend/internal/dashboard/handler"
	"expent-backend/internal/dashboard/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	svc := service.NewService()
	h := handler.NewHandler(svc)
	dashGroup := rg.Group("/dashboard")
	dashGroup.GET("", h.GetDashboard)
}
