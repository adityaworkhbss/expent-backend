package repository

import (
	"context"
	"expent-backend/internal/auth/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
	"time"

	"go.uber.org/zap"
)

type Repository interface {
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(email, name, id string) (*model.User, error)
	StoreRefreshToken(userID, token string) error
	GetRefreshToken(token string) (*model.RefreshToken, error)
	UpdateOnboardingStep(userID string, step int) error
	DeleteUserAndData(ctx context.Context, userID string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) GetUserByEmail(email string) (*model.User, error) {
	user, err := r.prisma.Prisma.User.FindFirst(
		db.User.Email.Equals(email),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		zap.S().Error("failed to fetch user", zap.Error(err))
		return nil, err
	}
	obs := 0
	if val, ok := user.OnboardingStep(); ok {
		obs = val
	}
	return &model.User{ID: user.ID, Email: user.Email, Name: func() string { v, _ := user.Name(); return v }(), OnboardingStep: obs}, nil
}

func (r *repoImpl) CreateUser(email, name, id string) (*model.User, error) {
	var params []db.UserSetParam
	if name != "" {
		params = append(params, db.User.Name.Set(name))
	}
	if id != "" {
		params = append(params, db.User.ID.Set(id))
	}
	usr, err := r.prisma.Prisma.User.CreateOne(
		db.User.Email.Set(email),
		params...,
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to create user", zap.Error(err))
		return nil, err
	}
	obs := 0
	if val, ok := usr.OnboardingStep(); ok {
		obs = val
	}
	return &model.User{ID: usr.ID, Email: usr.Email, Name: func() string { v, _ := usr.Name(); return v }(), OnboardingStep: obs}, nil
}

func (r *repoImpl) StoreRefreshToken(userID, token string) error {
	_, err := r.prisma.Prisma.RefreshToken.CreateOne(
		db.RefreshToken.Token.Set(token),
		db.RefreshToken.ExpiresAt.Set(time.Now().Add(7*24*time.Hour)),
		db.RefreshToken.User.Link(db.User.ID.Equals(userID)),
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to store refresh token", zap.Error(err))
	}
	return err
}

func (r *repoImpl) GetRefreshToken(token string) (*model.RefreshToken, error) {
	rt, err := r.prisma.Prisma.RefreshToken.FindFirst(
		db.RefreshToken.Token.Equals(token),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		zap.S().Error("failed to fetch refresh token", zap.Error(err))
		return nil, err
	}
	return &model.RefreshToken{ID: rt.ID, Token: rt.Token, ExpiresAt: rt.ExpiresAt, UserID: rt.UserID}, nil
}

func (r *repoImpl) UpdateOnboardingStep(userID string, step int) error {
	_, err := r.prisma.Prisma.User.FindUnique(
		db.User.ID.Equals(userID),
	).Update(
		db.User.OnboardingStep.Set(step),
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to update onboarding step", zap.Error(err))
	}
	return err
}

// DeleteUserAndData permanently deletes a user and everything they own, in one database transaction.
// The schema has no cascading deletes, so rows that reference others go first: transactions, EMIs and
// budgets point at categories/accounts, which point at the user.
func (r *repoImpl) DeleteUserAndData(ctx context.Context, userID string) error {
	client := r.prisma.Prisma
	err := client.Prisma.Transaction(
		client.Transaction.FindMany(db.Transaction.UserID.Equals(userID)).Delete().Tx(),
		client.Emi.FindMany(db.Emi.UserID.Equals(userID)).Delete().Tx(),
		client.Budget.FindMany(db.Budget.UserID.Equals(userID)).Delete().Tx(),
		client.Category.FindMany(db.Category.UserID.Equals(userID)).Delete().Tx(),
		client.Account.FindMany(db.Account.UserID.Equals(userID)).Delete().Tx(),
		client.UserCustomization.FindMany(db.UserCustomization.UserID.Equals(userID)).Delete().Tx(),
		client.DeviceToken.FindMany(db.DeviceToken.UserID.Equals(userID)).Delete().Tx(),
		client.Feedback.FindMany(db.Feedback.UserID.Equals(userID)).Delete().Tx(),
		client.RefreshToken.FindMany(db.RefreshToken.UserID.Equals(userID)).Delete().Tx(),
		client.User.FindMany(db.User.ID.Equals(userID)).Delete().Tx(),
	).Exec(ctx)
	if err != nil {
		zap.S().Error("failed to delete user and data", zap.String("userId", userID), zap.Error(err))
		return err
	}
	return nil
}
