package repository

import (
	"context"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/internal/transaction/model"
	"expent-backend/prisma/db"
	"time"
)

type Repository interface {
	ListTransactions(userID string, from, to *time.Time, page, limit *int) ([]model.Transaction, int, error)
	CreateTransaction(t model.Transaction) (*model.Transaction, error)
	GetTransactionByID(id string) (*model.Transaction, error)
	UpdateTransaction(id string, t model.Transaction) (*model.Transaction, error)
	DeleteTransaction(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func toModel(t *db.TransactionModel) model.Transaction {
	return model.Transaction{
		ID:         t.ID,
		UserID:     t.UserID,
		AccountID:  t.AccountID,
		CategoryID: t.CategoryID,
		Amount:     t.Amount,
		Type:       t.Type,
		Date:       t.Date.UTC().Format(time.RFC3339),
		Timestamp:  t.Date,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func (r *repoImpl) ListTransactions(userID string, from, to *time.Time, page, limit *int) ([]model.Transaction, int, error) {
	ctx := context.Background()

	// Base conditions
	var queryConditions []db.TransactionWhereParam
	queryConditions = append(queryConditions, db.Transaction.UserID.Equals(userID))

	if from != nil {
		queryConditions = append(queryConditions, db.Transaction.Date.Gte(*from))
	}
	if to != nil {
		queryConditions = append(queryConditions, db.Transaction.Date.Lte(*to))
	}

	// Count total
	allTs, err := r.prisma.Prisma.Transaction.FindMany(queryConditions...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := len(allTs)

	// Find many with ordering and pagination
	findQuery := r.prisma.Prisma.Transaction.FindMany(queryConditions...).OrderBy(
		db.Transaction.Date.Order(db.SortOrderDesc),
	)

	if page != nil && limit != nil {
		skip := (*page - 1) * (*limit)
		findQuery = findQuery.Skip(skip).Take(*limit)
	}

	ts, err := findQuery.Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.Transaction, 0, len(ts))
	for i := range ts {
		result = append(result, toModel(&ts[i]))
	}
	return result, total, nil
}

func (r *repoImpl) CreateTransaction(t model.Transaction) (*model.Transaction, error) {
	tx, err := r.prisma.Prisma.Transaction.CreateOne(
		db.Transaction.Amount.Set(t.Amount),
		db.Transaction.Type.Set(t.Type),
		db.Transaction.Date.Set(t.Timestamp),
		db.Transaction.User.Link(db.User.ID.Equals(t.UserID)),
		db.Transaction.Category.Link(db.Category.ID.Equals(t.CategoryID)),
		db.Transaction.Account.Link(db.Account.ID.Equals(t.AccountID)),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	m := toModel(tx)
	return &m, nil
}

func (r *repoImpl) GetTransactionByID(id string) (*model.Transaction, error) {
	t, err := r.prisma.Prisma.Transaction.FindFirst(
		db.Transaction.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	m := toModel(t)
	return &m, nil
}

func (r *repoImpl) UpdateTransaction(id string, t model.Transaction) (*model.Transaction, error) {
	tx, err := r.prisma.Prisma.Transaction.FindUnique(
		db.Transaction.ID.Equals(id),
	).Update(
		db.Transaction.Amount.Set(t.Amount),
		db.Transaction.Type.Set(t.Type),
		db.Transaction.Date.Set(t.Timestamp),
		db.Transaction.Category.Link(db.Category.ID.Equals(t.CategoryID)),
		db.Transaction.Account.Link(db.Account.ID.Equals(t.AccountID)),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	m := toModel(tx)
	return &m, nil
}

func (r *repoImpl) DeleteTransaction(id string) error {
	_, err := r.prisma.Prisma.Transaction.FindUnique(db.Transaction.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}
