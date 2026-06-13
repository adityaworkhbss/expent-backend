package handler

import (
	"encoding/json"
	"expent-backend/internal/emi/model"
	"expent-backend/internal/emi/service"
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

func (h *Handler) ListEmis(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	emis, err := h.svc.ListEmis(c.Request.Context(), userID.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "EMIs fetched", emis)
}

func (h *Handler) CreateEmi(c *gin.Context) {
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
	type EmiReq struct {
		Amount      float64 `json:"amount" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		CategoryID  *string `json:"categoryId"`
		AccountID   *string `json:"accountId"`
		Date        string  `json:"date" binding:"required"`
		Description *string `json:"description"`
	}

	var reqs []EmiReq

	if strings.HasPrefix(bodyStr, "[") {
		if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON array payload: "+err.Error())
			return
		}
	} else {
		var single EmiReq
		if err := json.Unmarshal(bodyBytes, &single); err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON object payload: "+err.Error())
			return
		}
		reqs = append(reqs, single)
	}

	var createdEmis []interface{}
	for _, req := range reqs {
		if req.Amount <= 0 || req.Type == "" || req.Date == "" {
			shared.ErrorResponse(c, http.StatusBadRequest, "amount, type, and date are required")
			return
		}

		t := parseDateHelper(req.Date)

		e := model.Emi{
			UserID:      userID.(string),
			Amount:      req.Amount,
			Type:        req.Type,
			CategoryID:  req.CategoryID,
			AccountID:   req.AccountID,
			Date:        t,
			DateStr:     t.Format("2006-01-02"),
			Description: req.Description,
		}

		res, err := h.svc.CreateEmi(c.Request.Context(), e)
		if err != nil {
			shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		createdEmis = append(createdEmis, res)
	}

	shared.SuccessResponse(c, "EMIs created", createdEmis)
}

func (h *Handler) UpdateEmi(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	emiID := c.Param("id")

	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		CategoryID  *string `json:"categoryId"`
		AccountID   *string `json:"accountId"`
		Date        string  `json:"date" binding:"required"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	t := parseDateHelper(req.Date)

	e := model.Emi{
		UserID:      userID.(string),
		Amount:      req.Amount,
		Type:        req.Type,
		CategoryID:  req.CategoryID,
		AccountID:   req.AccountID,
		Date:        t,
		DateStr:     t.Format("2006-01-02"),
		Description: req.Description,
	}

	res, err := h.svc.UpdateEmi(c.Request.Context(), userID.(string), emiID, e)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "EMI updated", res)
}

func (h *Handler) DeleteEmi(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	emiID := c.Param("id")

	if err := h.svc.DeleteEmi(c.Request.Context(), userID.(string), emiID); err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "EMI deleted", nil)
}

func parseDateHelper(dStr string) time.Time {
	t, err := time.Parse("2006-01-02", dStr)
	if err == nil {
		return t
	}
	t, err = time.Parse(time.RFC3339, dStr)
	if err == nil {
		return t
	}
	return time.Now()
}
