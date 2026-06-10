package repository

import (
	"context"
	"expent-backend/internal/transaction/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	ListTransactions(userID string) ([]model.Transaction, error)
	CreateTransaction(t model.Transaction) (*model.Transaction, error)
	GetTransactionByID(id string) (*model.Transaction, error)
	DeleteTransaction(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) ListTransactions(userID string) ([]model.Transaction, error) {
	ts, err := r.prisma.Prisma.Transaction.FindMany(
		db.Transaction.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var result []model.Transaction
	for _, t := range ts {
		result = append(result, model.Transaction{ID: t.ID, UserID: t.UserID, AccountID: t.AccountID, CategoryID: t.CategoryID, Amount: t.Amount, Type: t.Type, Timestamp: t.Date, Description: func() string { v, _ := t.Notes(); return v }()})
	}
	return result, nil
}

func (r *repoImpl) CreateTransaction(t model.Transaction) (*model.Transaction, error) {
	var optional []db.TransactionSetParam
	if t.Description != "" {
		optional = append(optional, db.Transaction.Notes.Set(t.Description))
	}

	tx, err := r.prisma.Prisma.Transaction.CreateOne(
		db.Transaction.Amount.Set(t.Amount),
		db.Transaction.Type.Set(t.Type),
		db.Transaction.Date.Set(t.Timestamp),
		db.Transaction.User.Link(db.User.ID.Equals(t.UserID)),
		db.Transaction.Category.Link(db.Category.ID.Equals(t.CategoryID)),
		db.Transaction.Account.Link(db.Account.ID.Equals(t.AccountID)),
		optional...
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &model.Transaction{ID: tx.ID, UserID: tx.UserID, AccountID: tx.AccountID, CategoryID: tx.CategoryID, Amount: tx.Amount, Type: tx.Type, Timestamp: tx.Date, Description: func() string { v, _ := tx.Notes(); return v }()}, nil
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
	return &model.Transaction{ID: t.ID, UserID: t.UserID, AccountID: t.AccountID, CategoryID: t.CategoryID, Amount: t.Amount, Type: t.Type, Timestamp: t.Date, Description: func() string { v, _ := t.Notes(); return v }()}, nil
}

func (r *repoImpl) DeleteTransaction(id string) error {
	_, err := r.prisma.Prisma.Transaction.FindUnique(db.Transaction.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}

