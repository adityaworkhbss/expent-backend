package handler

import (
	"context"
	"encoding/json"
	"expent-backend/internal/category/service"
	"expent-backend/internal/shared"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	bodyBytes, err := c.GetRawData()
	if err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Failed to read request body")
		return
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))
	var reqs []struct {
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Color string `json:"color"`
		Icon  string `json:"icon"`
	}

	if strings.HasPrefix(bodyStr, "[") {
		if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON array payload: "+err.Error())
			return
		}
	} else {
		var single struct {
			Name  string `json:"name" binding:"required"`
			Type  string `json:"type" binding:"required"`
			Color string `json:"color"`
			Icon  string `json:"icon"`
		}
		if err := json.Unmarshal(bodyBytes, &single); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON object payload: "+err.Error())
			return
		}
		if single.Name == "" || single.Type == "" {
			shared.ErrorResponse(c, http.StatusBadRequest, "name and type are required")
			return
		}
		reqs = append(reqs, single)
	}

	var createdCats []interface{}
	for _, req := range reqs {
		if req.Name == "" || req.Type == "" {
			shared.ErrorResponse(c, http.StatusBadRequest, "name and type are required")
			return
		}
		cat, err := h.svc.CreateCategory(c.Request.Context(), userID.(string), req.Name, req.Type, req.Color, req.Icon)
		if err != nil {
			shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		createdCats = append(createdCats, cat)
	}

	shared.SuccessResponse(c, "Categories created", createdCats)
}

// UpdateCategory PUT /categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	catID := c.Param("id")
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

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

	cat, err := h.svc.UpdateCategory(c.Request.Context(), userID.(string), catID, req.Name, req.Type, req.Color, req.Icon)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Category updated", cat)
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
