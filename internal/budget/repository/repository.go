package repository

import (
	"context"
	"time"
	"expent-backend/internal/budget/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	ListBudgets(userID string) ([]model.Budget, error)
	CreateBudget(budget model.Budget) (*model.Budget, error)
	GetBudgetByID(id string) (*model.Budget, error)
	DeleteBudget(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) ListBudgets(userID string) ([]model.Budget, error) {
	bs, err := r.prisma.Prisma.Budget.FindMany(
		db.Budget.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var result []model.Budget
	for _, b := range bs {
		result = append(result, model.Budget{ID: b.ID, UserID: b.UserID, CategoryID: b.CategoryID, Amount: b.LimitAmount, Period: b.PeriodType})
	}
	return result, nil
}

func (r *repoImpl) CreateBudget(budget model.Budget) (*model.Budget, error) {
	b, err := r.prisma.Prisma.Budget.CreateOne(
		db.Budget.PeriodType.Set(budget.Period),
		db.Budget.LimitAmount.Set(budget.Amount),
		db.Budget.StartDate.Set(time.Now()), // placeholder
		db.Budget.EndDate.Set(time.Now().Add(30*24*time.Hour)), // placeholder
		db.Budget.User.Link(db.User.ID.Equals(budget.UserID)),
		db.Budget.Category.Link(db.Category.ID.Equals(budget.CategoryID)),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &model.Budget{ID: b.ID, UserID: b.UserID, CategoryID: b.CategoryID, Amount: b.LimitAmount, Period: b.PeriodType}, nil
}

func (r *repoImpl) GetBudgetByID(id string) (*model.Budget, error) {
	b, err := r.prisma.Prisma.Budget.FindFirst(
		db.Budget.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Budget{ID: b.ID, UserID: b.UserID, CategoryID: b.CategoryID, Amount: b.LimitAmount, Period: b.PeriodType}, nil
}

func (r *repoImpl) DeleteBudget(id string) error {
	_, err := r.prisma.Prisma.Budget.FindUnique(db.Budget.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}

