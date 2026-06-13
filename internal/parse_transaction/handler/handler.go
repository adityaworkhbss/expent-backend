package handler

import (
	"net/http"

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
// It accepts a raw bank SMS or text and returns structured transaction data using Gemini AI.
func (h *Handler) ParseTransaction(c *gin.Context) {
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

	parsed, err := h.svc.ParseTransaction(c.Request.Context(), text)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	shared.SuccessResponse(c, "Transaction parsed successfully", parsed)
}
