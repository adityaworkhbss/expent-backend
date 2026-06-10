package repository

import (
	"context"
	"expent-backend/internal/category/model"
	"expent-backend/internal/infrastructure/prisma"
	"expent-backend/prisma/db"
)

type Repository interface {
	ListCategories(userID string) ([]model.Category, error)
	CreateCategory(cat model.Category) (*model.Category, error)
	GetCategoryByID(id string) (*model.Category, error)
	DeleteCategory(id string) error
}

type repoImpl struct {
	prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
	return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) ListCategories(userID string) ([]model.Category, error) {
	cats, err := r.prisma.Prisma.Category.FindMany(
		db.Category.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var result []model.Category
	for _, c := range cats {
		result = append(result, model.Category{ID: c.ID, UserID: c.UserID, Name: c.Name, Type: c.Type, Color: func() string { v, _ := c.Color(); return v }(), Icon: func() string { v, _ := c.Icon(); return v }()})
	}
	return result, nil
}

func (r *repoImpl) CreateCategory(cat model.Category) (*model.Category, error) {
	var optional []db.CategorySetParam
	if cat.Color != "" {
		optional = append(optional, db.Category.Color.Set(cat.Color))
	}
	if cat.Icon != "" {
		optional = append(optional, db.Category.Icon.Set(cat.Icon))
	}

	c, err := r.prisma.Prisma.Category.CreateOne(
		db.Category.Name.Set(cat.Name),
		db.Category.Type.Set(cat.Type),
		db.Category.User.Link(db.User.ID.Equals(cat.UserID)),
		optional...
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &model.Category{ID: c.ID, UserID: c.UserID, Name: c.Name, Type: c.Type, Color: func() string { v, _ := c.Color(); return v }(), Icon: func() string { v, _ := c.Icon(); return v }()}, nil
}

func (r *repoImpl) GetCategoryByID(id string) (*model.Category, error) {
	c, err := r.prisma.Prisma.Category.FindFirst(
		db.Category.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		if db.IsErrNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Category{ID: c.ID, UserID: c.UserID, Name: c.Name, Type: c.Type, Color: func() string { v, _ := c.Color(); return v }(), Icon: func() string { v, _ := c.Icon(); return v }()}, nil
}

func (r *repoImpl) DeleteCategory(id string) error {
	_, err := r.prisma.Prisma.Category.FindUnique(db.Category.ID.Equals(id)).Delete().Exec(context.Background())
	return err
}

