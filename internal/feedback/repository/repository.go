package repository

import (
	"context"

	"expent-backend/internal/feedback/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"

	"go.uber.org/zap"
)

type Repository interface {
	CreateFeedback(userID string, rating int, message string) (*model.Feedback, error)
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) CreateFeedback(userID string, rating int, message string) (*model.Feedback, error) {
	params := []db.FeedbackSetParam{}
	if message != "" {
		params = append(params, db.Feedback.Message.Set(message))
	}

	f, err := r.prisma.Prisma.Feedback.CreateOne(
		db.Feedback.Rating.Set(rating),
		db.Feedback.User.Link(db.User.ID.Equals(userID)),
		params...,
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to create feedback", zap.Error(err))
		return nil, err
	}

	msg, _ := f.Message()
	return &model.Feedback{
		ID:        f.ID,
		UserID:    f.UserID,
		Rating:    f.Rating,
		Message:   msg,
		CreatedAt: f.CreatedAt,
	}, nil
}
