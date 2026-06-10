package transaction

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
    "expent-backend/internal/transaction/service"
    "expent-backend/internal/transaction/model"
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
    txs, err := h.svc.ListTransactions(c.Request.Context(), userID.(string))
    if err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Transactions fetched", txs)
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
    tx, err := h.svc.CreateTransaction(c.Request.Context(), req)
    if err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Transaction created", tx)
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
