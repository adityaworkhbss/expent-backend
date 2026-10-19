package repository

import (
	"expent-backend/internal/device_token/model"
)

type DeviceTokenResponse struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

func ToDeviceTokenResponse(dt *model.DeviceToken) DeviceTokenResponse {
	return DeviceTokenResponse{
		Token:  dt.Token,
		UserID: dt.UserID,
	}
}