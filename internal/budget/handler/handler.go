package handler

import (
	"encoding/json"
	"expent-backend/internal/budget/model"
	"expent-backend/internal/budget/service"
	"expent-backend/internal/shared"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListBudgets(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	budgets, err := h.svc.ListBudgets(c.Request.Context(), userID.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Budgets fetched", budgets)
}

func parseDateHelper(dStr string, defaultTime time.Time) time.Time {
	if dStr == "" {
		return defaultTime
	}
	t, err := time.Parse("2006-01-02", dStr)
	if err == nil {
		return t
	}
	t, err = time.Parse(time.RFC3339, dStr)
	if err == nil {
		return t
	}
	return defaultTime
}

func (h *Handler) CreateBudget(c *gin.Context) {
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
	type BudgetReq struct {
		CategoryID  string  `json:"categoryId" binding:"required"`
		PeriodType  string  `json:"periodType" binding:"required"`
		LimitAmount float64 `json:"limitAmount" binding:"required"`
		StartDate   string  `json:"startDate"`
		EndDate     string  `json:"endDate"`
	}

	var reqs []BudgetReq

	if strings.HasPrefix(bodyStr, "[") {
		if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON array payload: "+err.Error())
			return
		}
	} else {
		var single BudgetReq
		if err := json.Unmarshal(bodyBytes, &single); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON object payload: "+err.Error())
			return
		}
		reqs = append(reqs, single)
	}

	var createdBudgets []interface{}
	for _, req := range reqs {
		if req.CategoryID == "" || req.PeriodType == "" || req.LimitAmount <= 0 {
			shared.ErrorResponse(c, http.StatusBadRequest, "categoryId, periodType, and positive limitAmount are required")
			return
		}

		startDate := parseDateHelper(req.StartDate, time.Now())
		endDate := parseDateHelper(req.EndDate, startDate.Add(30*24*time.Hour))

		b := model.Budget{
			UserID:     userID.(string),
			CategoryID: req.CategoryID,
			Period:     req.PeriodType,
			Amount:     req.LimitAmount,
			StartDate:  startDate,
			EndDate:    endDate,
		}

		res, err := h.svc.CreateBudget(c.Request.Context(), b)
		if err != nil {
			shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		createdBudgets = append(createdBudgets, res)
	}

	shared.SuccessResponse(c, "Budgets created", createdBudgets)
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	budgetID := c.Param("id")

	var req struct {
		CategoryID  string  `json:"categoryId" binding:"required"`
		PeriodType  string  `json:"periodType" binding:"required"`
		LimitAmount float64 `json:"limitAmount" binding:"required"`
		StartDate   string  `json:"startDate"`
		EndDate     string  `json:"endDate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if req.CategoryID == "" {
		shared.ErrorResponse(c, http.StatusBadRequest, "categoryId is required")
		return
	}

	startDate := parseDateHelper(req.StartDate, time.Now())
	endDate := parseDateHelper(req.EndDate, startDate.Add(30*24*time.Hour))

	b := model.Budget{
		UserID:     userID.(string),
		CategoryID: req.CategoryID,
		Period:     req.PeriodType,
		Amount:     req.LimitAmount,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	res, err := h.svc.UpdateBudget(c.Request.Context(), userID.(string), budgetID, b)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Budget updated", res)
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	budgetID := c.Param("id")

	if err := h.svc.DeleteBudget(c.Request.Context(), userID.(string), budgetID); err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Budget deleted", nil)
}
