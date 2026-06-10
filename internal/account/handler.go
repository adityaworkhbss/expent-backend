package account

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
    "expent-backend/internal/account/service"
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
    var req struct {
        Name string `json:"name" binding:"required"`
        Type string `json:"type" binding:"required"`
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
    acc, err := h.svc.CreateAccount(context.Background(), userID.(string), req.Name, req.Type)
    if err != nil {
        shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    shared.SuccessResponse(c, "Account created", acc)
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
