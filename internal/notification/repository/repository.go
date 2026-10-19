package repository

import (
	"context"

	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/notification/model"
	"expent-backend/prisma/db"

	"go.uber.org/zap"
)

type Repository interface {
	CreateNotification(title, body, createdBy string, sentCount, failedCount int) (*model.Notification, error)
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) CreateNotification(title, body, createdBy string, sentCount, failedCount int) (*model.Notification, error) {
	n, err := r.prisma.Prisma.Notification.CreateOne(
		db.Notification.Title.Set(title),
		db.Notification.Body.Set(body),
		db.Notification.CreatedBy.Set(createdBy),
		db.Notification.SentCount.Set(sentCount),
		db.Notification.FailedCount.Set(failedCount),
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to create notification", zap.Error(err))
		return nil, err
	}
	return &model.Notification{
		ID:          n.ID,
		Title:       n.Title,
		Body:        n.Body,
		CreatedBy:   n.CreatedBy,
		SentCount:   n.SentCount,
		FailedCount: n.FailedCount,
		CreatedAt:   n.CreatedAt,
	}, nil
}
