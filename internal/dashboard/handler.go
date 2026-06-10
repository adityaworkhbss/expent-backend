package dashboard

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
)

type Handler struct {}

func NewHandler() *Handler { return &Handler{} }

// GetDashboard GET /dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
    // Placeholder response – real implementation would aggregate data.
    shared.SuccessResponse(c, "Dashboard data", gin.H{"message": "Dashboard endpoint works"})
}
