package handler

import (
	"context"
	"encoding/json"
	"expent-backend/internal/account/service"
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

// ListAccounts GET /accounts
func (h *Handler) ListAccounts(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	accounts, err := h.svc.ListAccounts(context.Background(), userID.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Accounts fetched", accounts)
}

// CreateAccount POST /accounts
func (h *Handler) CreateAccount(c *gin.Context) {
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
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if strings.HasPrefix(bodyStr, "[") {
		if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON array payload: "+err.Error())
			return
		}
	} else {
		var single struct {
			Name string `json:"name" binding:"required"`
			Type string `json:"type" binding:"required"`
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

	var createdAccs []interface{}
	for _, req := range reqs {
		if req.Name == "" || req.Type == "" {
			shared.ErrorResponse(c, http.StatusBadRequest, "name and type are required")
			return
		}
		acc, err := h.svc.CreateAccount(c.Request.Context(), userID.(string), req.Name, req.Type)
		if err != nil {
			shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		createdAccs = append(createdAccs, acc)
	}

	shared.SuccessResponse(c, "Accounts created", createdAccs)
}

// UpdateAccount PUT /accounts/:id
func (h *Handler) UpdateAccount(c *gin.Context) {
	accID := c.Param("id")
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	acc, err := h.svc.UpdateAccount(c.Request.Context(), userID.(string), accID, req.Name, req.Type)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Account updated", acc)
}

// DeleteAccount DELETE /accounts/:id
func (h *Handler) DeleteAccount(c *gin.Context) {
	accID := c.Param("id")
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	if err := h.svc.DeleteAccount(context.Background(), userID.(string), accID); err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Account deleted", nil)
}
