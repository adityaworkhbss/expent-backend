package repository

import (
	"context"

	"expent-backend/internal/device_token/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"

	"go.uber.org/zap"
)

type Repository interface {
	RegisterDeviceToken(userID, token string) (*model.DeviceToken, error)
	GetDeviceTokenByToken(token string) (*model.DeviceToken, error)
	GetDeviceTokensByUserID(userID string) ([]*model.DeviceToken, error)
	GetAllDeviceTokens() ([]*model.DeviceToken, error)
	DeleteDeviceTokensByTokens(tokens []string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) RegisterDeviceToken(userID, token string) (*model.DeviceToken, error) {
	existing, err := r.prisma.Prisma.DeviceToken.FindFirst(
		db.DeviceToken.Token.Equals(token),
	).Exec(context.Background())

	if err != nil && !db.IsErrNotFound(err) {
		zap.S().Error("failed to check existing device token", zap.Error(err))
		return nil, err
	}

	var dt *db.DeviceTokenModel
	if existing != nil {
		dt, err = r.prisma.Prisma.DeviceToken.FindUnique(
			db.DeviceToken.Token.Equals(token),
		).Update(
			db.DeviceToken.UserID.Set(userID),
		).Exec(context.Background())
	} else {
		dt, err = r.prisma.Prisma.DeviceToken.CreateOne(
			db.DeviceToken.Token.Set(token),
			db.DeviceToken.User.Link(db.User.ID.Equals(userID)),
		).Exec(context.Background())
	}

	if err != nil {
		zap.S().Error("failed to register device token", zap.Error(err))
		return nil, err
	}

	return &model.DeviceToken{
		ID:        dt.ID,
		Token:     dt.Token,
		UserID:    dt.UserID,
		CreatedAt: dt.CreatedAt,
		UpdatedAt: dt.UpdatedAt,
	}, nil
}

func (r *repoImpl) GetDeviceTokenByToken(token string) (*model.DeviceToken, error) {
	dt, err := r.prisma.Prisma.DeviceToken.FindFirst(
		db.DeviceToken.Token.Equals(token),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		zap.S().Error("failed to fetch device token", zap.Error(err))
		return nil, err
	}
	return &model.DeviceToken{
		ID:        dt.ID,
		Token:     dt.Token,
		UserID:    dt.UserID,
		CreatedAt: dt.CreatedAt,
		UpdatedAt: dt.UpdatedAt,
	}, nil
}

func (r *repoImpl) GetDeviceTokensByUserID(userID string) ([]*model.DeviceToken, error) {
	dts, err := r.prisma.Prisma.DeviceToken.FindMany(
		db.DeviceToken.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to fetch device tokens", zap.Error(err))
		return nil, err
	}
	result := make([]*model.DeviceToken, len(dts))
	for i, dt := range dts {
		result[i] = &model.DeviceToken{
			ID:        dt.ID,
			Token:     dt.Token,
			UserID:    dt.UserID,
			CreatedAt: dt.CreatedAt,
			UpdatedAt: dt.UpdatedAt,
		}
	}
	return result, nil
}

func (r *repoImpl) GetAllDeviceTokens() ([]*model.DeviceToken, error) {
	dts, err := r.prisma.Prisma.DeviceToken.FindMany().Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to fetch all device tokens", zap.Error(err))
		return nil, err
	}
	result := make([]*model.DeviceToken, len(dts))
	for i, dt := range dts {
		result[i] = &model.DeviceToken{
			ID:        dt.ID,
			Token:     dt.Token,
			UserID:    dt.UserID,
			CreatedAt: dt.CreatedAt,
			UpdatedAt: dt.UpdatedAt,
		}
	}
	return result, nil
}

func (r *repoImpl) DeleteDeviceTokensByTokens(tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	_, err := r.prisma.Prisma.DeviceToken.FindMany(
		db.DeviceToken.Token.In(tokens),
	).Delete().Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to delete stale device tokens", zap.Error(err))
		return err
	}
	return nil
}