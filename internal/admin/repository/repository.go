package repository

import (
	"context"

	"expent-backend/internal/admin/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"

	"go.uber.org/zap"
)

type Repository interface {
	GetAdminByEmail(email string) (*model.Admin, error)
	CreateAdmin(email, passwordHash, name string) (*model.Admin, error)
	CountAdmins() (int, error)
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) GetAdminByEmail(email string) (*model.Admin, error) {
	admin, err := r.prisma.Prisma.Admin.FindFirst(
		db.Admin.Email.Equals(email),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		zap.S().Error("failed to fetch admin", zap.Error(err))
		return nil, err
	}
	name, _ := admin.Name()
	return &model.Admin{
		ID:           admin.ID,
		Email:        admin.Email,
		PasswordHash: admin.PasswordHash,
		Name:         name,
		CreatedAt:    admin.CreatedAt,
		UpdatedAt:    admin.UpdatedAt,
	}, nil
}

func (r *repoImpl) CountAdmins() (int, error) {
	admins, err := r.prisma.Prisma.Admin.FindMany().Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to count admins", zap.Error(err))
		return 0, err
	}
	return len(admins), nil
}

func (r *repoImpl) CreateAdmin(email, passwordHash, name string) (*model.Admin, error) {
	var params []db.AdminSetParam
	if name != "" {
		params = append(params, db.Admin.Name.Set(name))
	}
	admin, err := r.prisma.Prisma.Admin.CreateOne(
		db.Admin.Email.Set(email),
		db.Admin.PasswordHash.Set(passwordHash),
		params...,
	).Exec(context.Background())
	if err != nil {
		zap.S().Error("failed to create admin", zap.Error(err))
		return nil, err
	}
	createdName, _ := admin.Name()
	return &model.Admin{
		ID:           admin.ID,
		Email:        admin.Email,
		PasswordHash: admin.PasswordHash,
		Name:         createdName,
		CreatedAt:    admin.CreatedAt,
		UpdatedAt:    admin.UpdatedAt,
	}, nil
}
