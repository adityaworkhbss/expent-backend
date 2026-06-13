package repository

import (
	"context"
	"expent-backend/internal/account/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	ListAccounts(userID string) ([]model.Account, error)
	CreateAccount(acc model.Account) (*model.Account, error)
	GetAccountByID(id string) (*model.Account, error)
	UpdateAccount(id string, acc model.Account) (*model.Account, error)
	DeleteAccount(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) ListAccounts(userID string) ([]model.Account, error) {
	accs, err := r.prisma.Prisma.Account.FindMany(
		db.Account.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var result []model.Account
	for _, a := range accs {
		result = append(result, model.Account{ID: a.ID, UserID: a.UserID, Name: a.Name, Type: a.Type})
	}
	return result, nil
}

func (r *repoImpl) CreateAccount(acc model.Account) (*model.Account, error) {
	a, err := r.prisma.Prisma.Account.CreateOne(
		db.Account.Name.Set(acc.Name),
		db.Account.Type.Set(acc.Type),
		db.Account.User.Link(db.User.ID.Equals(acc.UserID)),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &model.Account{ID: a.ID, UserID: a.UserID, Name: a.Name, Type: a.Type}, nil
}

func (r *repoImpl) GetAccountByID(id string) (*model.Account, error) {
	a, err := r.prisma.Prisma.Account.FindFirst(
		db.Account.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Account{ID: a.ID, UserID: a.UserID, Name: a.Name, Type: a.Type}, nil
}

func (r *repoImpl) UpdateAccount(id string, acc model.Account) (*model.Account, error) {
	a, err := r.prisma.Prisma.Account.FindUnique(
		db.Account.ID.Equals(id),
	).Update(
		db.Account.Name.Set(acc.Name),
		db.Account.Type.Set(acc.Type),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &model.Account{ID: a.ID, UserID: a.UserID, Name: a.Name, Type: a.Type}, nil
}

func (r *repoImpl) DeleteAccount(id string) error {
	_, err := r.prisma.Prisma.Account.FindUnique(db.Account.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}

