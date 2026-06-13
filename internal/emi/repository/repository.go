package repository

import (
	"context"
	"expent-backend/internal/emi/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	ListEmis(userID string) ([]model.Emi, error)
	CreateEmi(emi model.Emi) (*model.Emi, error)
	GetEmiByID(id string) (*model.Emi, error)
	UpdateEmi(id string, emi model.Emi) (*model.Emi, error)
	DeleteEmi(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) mapEmiStruct(e *db.EmiModel) *model.Emi {
	catID, _ := e.CategoryID()
	accID, _ := e.AccountID()
	desc, _ := e.Description()

	var catIDPtr *string
	if catID != "" {
		catIDPtr = &catID
	}
	var accIDPtr *string
	if accID != "" {
		accIDPtr = &accID
	}
	var descPtr *string
	if desc != "" {
		descPtr = &desc
	}

	return &model.Emi{
		ID:          e.ID,
		UserID:      e.UserID,
		Amount:      e.Amount,
		Type:        e.Type,
		CategoryID:  catIDPtr,
		AccountID:   accIDPtr,
		Date:        e.Date,
		DateStr:     e.Date.Format("2006-01-02"),
		Description: descPtr,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (r *repoImpl) ListEmis(userID string) ([]model.Emi, error) {
	ctx := context.Background()
	es, err := r.prisma.Prisma.Emi.FindMany(
		db.Emi.UserID.Equals(userID),
	).OrderBy(
		db.Emi.Date.Order(db.SortOrderDesc),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.Emi
	for _, e := range es {
		result = append(result, *r.mapEmiStruct(&e))
	}
	return result, nil
}

func (r *repoImpl) CreateEmi(emi model.Emi) (*model.Emi, error) {
	ctx := context.Background()

	var optionalParams []db.EmiSetParam
	if emi.CategoryID != nil && *emi.CategoryID != "" {
		optionalParams = append(optionalParams, db.Emi.Category.Link(db.Category.ID.Equals(*emi.CategoryID)))
	}
	if emi.AccountID != nil && *emi.AccountID != "" {
		optionalParams = append(optionalParams, db.Emi.Account.Link(db.Account.ID.Equals(*emi.AccountID)))
	}
	if emi.Description != nil {
		optionalParams = append(optionalParams, db.Emi.Description.Set(*emi.Description))
	}

	tx, err := r.prisma.Prisma.Emi.CreateOne(
		db.Emi.Amount.Set(emi.Amount),
		db.Emi.Type.Set(emi.Type),
		db.Emi.Date.Set(emi.Date),
		db.Emi.User.Link(db.User.ID.Equals(emi.UserID)),
		optionalParams...,
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapEmiStruct(tx), nil
}

func (r *repoImpl) GetEmiByID(id string) (*model.Emi, error) {
	ctx := context.Background()
	e, err := r.prisma.Prisma.Emi.FindFirst(
		db.Emi.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.mapEmiStruct(e), nil
}

func (r *repoImpl) UpdateEmi(id string, emi model.Emi) (*model.Emi, error) {
	ctx := context.Background()

	var updateParams []db.EmiSetParam
	updateParams = append(updateParams, db.Emi.Amount.Set(emi.Amount))
	updateParams = append(updateParams, db.Emi.Type.Set(emi.Type))
	updateParams = append(updateParams, db.Emi.Date.Set(emi.Date))

	if emi.CategoryID != nil && *emi.CategoryID != "" {
		updateParams = append(updateParams, db.Emi.Category.Link(db.Category.ID.Equals(*emi.CategoryID)))
	} else {
		updateParams = append(updateParams, db.Emi.Category.Unlink())
	}

	if emi.AccountID != nil && *emi.AccountID != "" {
		updateParams = append(updateParams, db.Emi.Account.Link(db.Account.ID.Equals(*emi.AccountID)))
	} else {
		updateParams = append(updateParams, db.Emi.Account.Unlink())
	}

	if emi.Description != nil {
		updateParams = append(updateParams, db.Emi.Description.Set(*emi.Description))
	} else {
		updateParams = append(updateParams, db.Emi.Description.SetOptional(nil))
	}

	tx, err := r.prisma.Prisma.Emi.FindUnique(
		db.Emi.ID.Equals(id),
	).Update(
		updateParams...,
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapEmiStruct(tx), nil
}

func (r *repoImpl) DeleteEmi(id string) error {
	ctx := context.Background()
	_, err := r.prisma.Prisma.Emi.FindUnique(db.Emi.ID.Equals(id)).Delete().Exec(ctx)
	return err
}
