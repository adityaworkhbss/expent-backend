package handler

import (
	"encoding/json"
	"errors"
	"expent-backend/internal/emi/model"
	"expent-backend/internal/emi/service"
	"expent-backend/internal/shared"
	"fmt"
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

type emiReq struct {
	Amount      shared.FlexFloat64 `json:"amount"`
	Type        string             `json:"type"`
	Name        string             `json:"name"`
	StartDate   string             `json:"startDate"`
	Date        string             `json:"date"`
	Tenure      *shared.FlexInt    `json:"tenure"`
	MonthsPaid  *shared.FlexInt    `json:"monthsPaid"`
	CategoryID  *string            `json:"categoryId"`
	AccountID   *string            `json:"accountId"`
	Description *string            `json:"description"`
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

	reqs, err := decodeEmiReqs(c)
	if err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var createdEmis []interface{}
	for _, req := range reqs {
		e, err := req.toModel(userID.(string))
		if err != nil {
			shared.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
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

	bodyBytes, err := c.GetRawData()
	if err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var req emiReq
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid payload: "+err.Error())
		return
	}

	e, err := req.toModel(userID.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
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

func decodeEmiReqs(c *gin.Context) ([]emiReq, error) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		return nil, errors.New("Failed to read request body")
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))
	var reqs []emiReq

	if strings.HasPrefix(bodyStr, "[") {
		if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
			return nil, fmt.Errorf("Invalid JSON array payload: %w", err)
		}
		return reqs, nil
	}

	var single emiReq
	if err := json.Unmarshal(bodyBytes, &single); err != nil {
		return nil, fmt.Errorf("Invalid JSON object payload: %w", err)
	}
	return []emiReq{single}, nil
}

func (req emiReq) toModel(userID string) (model.Emi, error) {
	amount := req.Amount.Float64()
	if amount <= 0 {
		return model.Emi{}, errRequired("amount must be a positive number")
	}
	if strings.TrimSpace(req.Type) == "" {
		return model.Emi{}, errRequired("type is required")
	}

	dateStr := strings.TrimSpace(req.Date)
	if dateStr == "" {
		dateStr = strings.TrimSpace(req.StartDate)
	}
	t, err := shared.ParseFlexibleDate(dateStr)
	if err != nil {
		return model.Emi{}, err
	}

	name := strings.TrimSpace(req.Name)
	var namePtr *string
	if name != "" {
		namePtr = &name
	}

	var descPtr *string
	if req.Description != nil && strings.TrimSpace(*req.Description) != "" {
		d := strings.TrimSpace(*req.Description)
		descPtr = &d
	} else if namePtr != nil {
		descPtr = namePtr
	}

	dateFormatted := t.Format("2006-01-02")
	return model.Emi{
		UserID:      userID,
		Amount:      amount,
		Type:        req.Type,
		CategoryID:  emptyToNil(req.CategoryID),
		AccountID:   emptyToNil(req.AccountID),
		Date:        t,
		DateStr:     dateFormatted,
		StartDate:   dateFormatted,
		Name:        namePtr,
		Description: descPtr,
		Tenure:      shared.FlexIntPtr(req.Tenure),
		MonthsPaid:  shared.FlexIntPtr(req.MonthsPaid),
	}, nil
}

type simpleError string

func (e simpleError) Error() string { return string(e) }

func errRequired(msg string) error { return simpleError(msg) }

func emptyToNil(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}
