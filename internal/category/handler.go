package category

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
    "expent-backend/internal/category/service"
)

type Handler struct {
    svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
    return &Handler{svc: svc}
}

// ListCategories GET /categories
func (h *Handler) ListCategories(c *gin.Context) {
    userID, ok := c.Get("userId")
    if !ok {
        shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
    categories, err := h.svc.ListCategories(context.Background(), userID.(string))
    if err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Categories fetched", categories)
}

// CreateCategory POST /categories
func (h *Handler) CreateCategory(c *gin.Context) {
    var req struct {
        Name  string `json:"name" binding:"required"`
        Type  string `json:"type" binding:"required"`
        Color string `json:"color"`
        Icon  string `json:"icon"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
        return
    }
    userID, ok := c.Get("userId")
    if !ok {
        shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
    cat, err := h.svc.CreateCategory(context.Background(), userID.(string), req.Name, req.Type, req.Color, req.Icon)
    if err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Category created", cat)
}

// DeleteCategory DELETE /categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
    catID := c.Param("id")
    userID, ok := c.Get("userId")
    if !ok {
        shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
    if err := h.svc.DeleteCategory(context.Background(), userID.(string), catID); err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Category deleted", nil)
}
