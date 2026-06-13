package repository

import (
	"context"
	"expent-backend/internal/user_customization/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	GetCustomization(userID string) (*model.UserCustomization, error)
	UpdateCustomization(userID string, c model.UserCustomization) (*model.UserCustomization, error)
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) GetCustomization(userID string) (*model.UserCustomization, error) {
	ctx := context.Background()
	c, err := r.prisma.Prisma.UserCustomization.FindUnique(
		db.UserCustomization.UserID.Equals(userID),
	).Exec(ctx)
	if err != nil {
		if db.IsErrNotFound(err) {
			// Create default customization record
			newC, err := r.prisma.Prisma.UserCustomization.CreateOne(
				db.UserCustomization.User.Link(db.User.ID.Equals(userID)),
				db.UserCustomization.Currency.Set("INR"),
				db.UserCustomization.Theme.Set("dark"),
			).Exec(ctx)
			if err != nil {
				return nil, err
			}
			return &model.UserCustomization{
				ID:       newC.ID,
				UserID:   newC.UserID,
				Currency: func() string { v, _ := newC.Currency(); return v }(),
				Theme:    func() string { v, _ := newC.Theme(); return v }(),
			}, nil
		}
		return nil, err
	}
	return &model.UserCustomization{
		ID:       c.ID,
		UserID:   c.UserID,
		Currency: func() string { v, _ := c.Currency(); return v }(),
		Theme:    func() string { v, _ := c.Theme(); return v }(),
	}, nil
}

func (r *repoImpl) UpdateCustomization(userID string, c model.UserCustomization) (*model.UserCustomization, error) {
	ctx := context.Background()
	// Check if customization exists first
	_, err := r.GetCustomization(userID)
	if err != nil {
		return nil, err
	}

	updated, err := r.prisma.Prisma.UserCustomization.FindUnique(
		db.UserCustomization.UserID.Equals(userID),
	).Update(
		db.UserCustomization.Currency.Set(c.Currency),
		db.UserCustomization.Theme.Set(c.Theme),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &model.UserCustomization{
		ID:       updated.ID,
		UserID:   updated.UserID,
		Currency: func() string { v, _ := updated.Currency(); return v }(),
		Theme:    func() string { v, _ := updated.Theme(); return v }(),
	}, nil
}
