package dashboard

import (
    "github.com/gin-gonic/gin"
    "expent-backend/internal/dashboard/handler"
    "expent-backend/internal/dashboard/service"
)

func RegisterRoutes(rg *gin.RouterGroup) {
    svc := service.NewService()
    h := handler.NewHandler(svc)
    dashGroup := rg.Group("/dashboard")
    dashGroup.GET("", h.GetDashboard)
}
