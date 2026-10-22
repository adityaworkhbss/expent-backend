package app_config

import (
	"expent-backend/configs"
	"expent-backend/internal/app_config/handler"
	"expent-backend/internal/app_config/service"
	"expent-backend/internal/infrastructure/cloudflare"
	"expent-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the admin-only app config (remote text + server URL) routes.
// r should NOT already require user auth — AdminAuth is applied here.
func RegisterRoutes(r *gin.RouterGroup) {
	cfg := configs.AppConfig
	kv := cloudflare.NewKVClient(cfg.CF_ACCOUNT_ID, cfg.CF_KV_NAMESPACE_ID, cfg.CF_API_TOKEN)
	svc := service.NewService(kv, cfg.CF_KV_CONFIG_KEY, cfg.CONFIG_PUBLIC_URL, cfg.CONFIG_API_KEY)
	h := handler.NewHandler(svc)

	group := r.Group("/admin/app-config")
	group.Use(middleware.AdminAuth())
	group.GET("", h.Get)
	group.PUT("", h.Update)
}
