package handler

import (
	"expent-backend/internal/shared"
	"expent-backend/internal/transaction/model"
	"expent-backend/internal/transaction/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// ListTransactions GET /transactions
func (h *Handler) ListTransactions(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var from *time.Time
	var to *time.Time

	if fromStr != "" {
		if f, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &f
		} else if f, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &f
		}
	}

	if toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			to = &t
		} else if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	txs, total, err := h.svc.ListTransactions(c.Request.Context(), userID.(string), from, to, &page, &limit)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Format response to match PaginatedTransactionsResponseDto
	response := gin.H{
		"items": txs,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	shared.SuccessResponse(c, "Transactions fetched", response)
}

// CreateTransaction POST /transactions
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req model.Transaction
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	req.UserID = userID.(string)
	req.Timestamp = parseTransactionDate(req.Date)

	tx, err := h.svc.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Transaction created", tx)
}

// UpdateTransaction PUT /transactions/:id
func (h *Handler) UpdateTransaction(c *gin.Context) {
	txID := c.Param("id")
	var req model.Transaction
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	req.ID = txID
	req.UserID = userID.(string)
	req.Timestamp = parseTransactionDate(req.Date)

	tx, err := h.svc.UpdateTransaction(c.Request.Context(), userID.(string), req)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Transaction updated", tx)
}

// DeleteTransaction DELETE /transactions/:id
func (h *Handler) DeleteTransaction(c *gin.Context) {
	txID := c.Param("id")
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	if err := h.svc.DeleteTransaction(c.Request.Context(), userID.(string), txID); err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Transaction deleted", nil)
}

// parseTransactionDate accepts an RFC 3339 timestamp (fractional seconds allowed) or yyyy-MM-dd;
// anything else falls back to now.
func parseTransactionDate(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC()
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Now().UTC()
}
