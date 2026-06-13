package repository

import (
	"context"
	"time"
	"expent-backend/internal/transaction/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
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

	var result []model.Transaction
	for _, t := range ts {
		notes, _ := t.Notes()
		result = append(result, model.Transaction{
			ID:          t.ID,
			UserID:      t.UserID,
			AccountID:   t.AccountID,
			CategoryID:  t.CategoryID,
			Amount:      t.Amount,
			Type:        t.Type,
			Timestamp:   t.Date,
			Date:        t.Date.Format("2006-01-02"),
			Description: notes,
			Notes:       notes,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	return result, total, nil
}

func (r *repoImpl) CreateTransaction(t model.Transaction) (*model.Transaction, error) {
	var optional []db.TransactionSetParam
	notesVal := t.Description
	if notesVal == "" {
		notesVal = t.Notes
	}
	if notesVal != "" {
		optional = append(optional, db.Transaction.Notes.Set(notesVal))
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
	notes, _ := tx.Notes()
	return &model.Transaction{
		ID:          tx.ID,
		UserID:      tx.UserID,
		AccountID:   tx.AccountID,
		CategoryID:  tx.CategoryID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Timestamp:   tx.Date,
		Date:        tx.Date.Format("2006-01-02"),
		Description: notes,
		Notes:       notes,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}, nil
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
	notes, _ := t.Notes()
	return &model.Transaction{
		ID:          t.ID,
		UserID:      t.UserID,
		AccountID:   t.AccountID,
		CategoryID:  t.CategoryID,
		Amount:      t.Amount,
		Type:        t.Type,
		Timestamp:   t.Date,
		Date:        t.Date.Format("2006-01-02"),
		Description: notes,
		Notes:       notes,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}

func (r *repoImpl) UpdateTransaction(id string, t model.Transaction) (*model.Transaction, error) {
	var updateParams []db.TransactionSetParam
	updateParams = append(updateParams, db.Transaction.Amount.Set(t.Amount))
	updateParams = append(updateParams, db.Transaction.Type.Set(t.Type))
	updateParams = append(updateParams, db.Transaction.Date.Set(t.Timestamp))
	updateParams = append(updateParams, db.Transaction.Category.Link(db.Category.ID.Equals(t.CategoryID)))
	updateParams = append(updateParams, db.Transaction.Account.Link(db.Account.ID.Equals(t.AccountID)))

	notesVal := t.Description
	if notesVal == "" {
		notesVal = t.Notes
	}
	if notesVal != "" {
		updateParams = append(updateParams, db.Transaction.Notes.Set(notesVal))
	}

	tx, err := r.prisma.Prisma.Transaction.FindUnique(
		db.Transaction.ID.Equals(id),
	).Update(
		updateParams...,
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	notes, _ := tx.Notes()
	return &model.Transaction{
		ID:          tx.ID,
		UserID:      tx.UserID,
		AccountID:   tx.AccountID,
		CategoryID:  tx.CategoryID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Timestamp:   tx.Date,
		Date:        tx.Date.Format("2006-01-02"),
		Description: notes,
		Notes:       notes,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}, nil
}

func (r *repoImpl) DeleteTransaction(id string) error {
	_, err := r.prisma.Prisma.Transaction.FindUnique(db.Transaction.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}

