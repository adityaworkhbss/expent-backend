package handler

import (
	"net/http"
	"time"

	"expent-backend/internal/parse_transaction/model"
	"expent-backend/internal/parse_transaction/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// ParseTransaction handles POST /parse-transaction
func (h *Handler) ParseTransaction(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req model.ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	text := req.Text
	if text == "" {
		text = req.RawText
	}
	if text == "" {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload: text or rawText is required")
		return
	}

	today := req.Today
	if today == "" {
		today = time.Now().Format("2006-01-02")
	}

	parsed, err := h.svc.ParseTransaction(c.Request.Context(), userID.(string), text, today, req.Categories, req.Accounts)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                 true,
		"message":                 "Transaction parsed successfully",
		"requires_user_input":     parsed.RequiresUserInput,
		"missing_required_fields": parsed.MissingRequiredFields,
		"data":                    parsed.Data,
	})
}
