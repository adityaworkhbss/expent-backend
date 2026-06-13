package repository

import (
	"context"
	"time"

	"expent-backend/internal/dashboard/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

// Repository defines the data-access methods required by the dashboard service.
type Repository interface {
	GetMonthlyTransactions(userID string, from, to time.Time) ([]db.TransactionModel, error)
	GetRecentTransactions(userID string, limit int) ([]db.TransactionModel, error)
	GetBudgets(userID string) ([]db.BudgetModel, error)
	GetTransactionsByCategoryAndDateRange(userID, categoryID string, from, to time.Time) ([]db.TransactionModel, error)
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

// NewRepository creates a new dashboard repository backed by Prisma.
func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) GetMonthlyTransactions(userID string, from, to time.Time) ([]db.TransactionModel, error) {
	ctx := context.Background()

	ts, err := r.prisma.Prisma.Transaction.FindMany(
		db.Transaction.UserID.Equals(userID),
		db.Transaction.Date.Gte(from),
		db.Transaction.Date.Lte(to),
	).With(
		db.Transaction.Category.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *repoImpl) GetRecentTransactions(userID string, limit int) ([]db.TransactionModel, error) {
	ctx := context.Background()

	ts, err := r.prisma.Prisma.Transaction.FindMany(
		db.Transaction.UserID.Equals(userID),
	).OrderBy(
		db.Transaction.Date.Order(db.SortOrderDesc),
	).Take(limit).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *repoImpl) GetBudgets(userID string) ([]db.BudgetModel, error) {
	ctx := context.Background()

	bs, err := r.prisma.Prisma.Budget.FindMany(
		db.Budget.UserID.Equals(userID),
	).With(
		db.Budget.Category.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (r *repoImpl) GetTransactionsByCategoryAndDateRange(userID, categoryID string, from, to time.Time) ([]db.TransactionModel, error) {
	ctx := context.Background()

	ts, err := r.prisma.Prisma.Transaction.FindMany(
		db.Transaction.UserID.Equals(userID),
		db.Transaction.CategoryID.Equals(categoryID),
		db.Transaction.Type.Equals("expense"),
		db.Transaction.Date.Gte(from),
		db.Transaction.Date.Lte(to),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// MapRecentTransactions converts Prisma transaction models to dashboard model structs.
func MapRecentTransactions(ts []db.TransactionModel) []model.RecentTransaction {
	var result []model.RecentTransaction
	for _, t := range ts {
		notes, _ := t.Notes()
		result = append(result, model.RecentTransaction{
			ID:         t.ID,
			Amount:     t.Amount,
			Type:       t.Type,
			Notes:      notes,
			Date:       t.Date.Format("2006-01-02"),
			CategoryID: t.CategoryID,
			AccountID:  t.AccountID,
			CreatedAt:  t.CreatedAt,
		})
	}
	return result
}
