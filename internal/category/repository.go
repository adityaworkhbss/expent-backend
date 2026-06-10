package category

import (
    "context"
    "expent-backend/internal/category/model"
    "expent-backend/internal/infrastructure/prisma"
    "github.com/steebchen/prisma-client-go/runtime"
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
        prisma.Category.UserId.Equals(userID),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    var result []model.Category
    for _, c := range cats {
        result = append(result, model.Category{ID: c.ID, UserID: c.UserId, Name: c.Name, Type: c.Type, Color: c.Color, Icon: c.Icon})
    }
    return result, nil
}

func (r *repoImpl) CreateCategory(cat model.Category) (*model.Category, error) {
    c, err := r.prisma.Prisma.Category.CreateOne(
        prisma.Category.UserId.Set(cat.UserID),
        prisma.Category.Name.Set(cat.Name),
        prisma.Category.Type.Set(cat.Type),
        prisma.Category.Color.Set(cat.Color),
        prisma.Category.Icon.Set(cat.Icon),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    return &model.Category{ID: c.ID, UserID: c.UserId, Name: c.Name, Type: c.Type, Color: c.Color, Icon: c.Icon}, nil
}

func (r *repoImpl) GetCategoryByID(id string) (*model.Category, error) {
    c, err := r.prisma.Prisma.Category.FindFirst(
        prisma.Category.ID.Equals(id),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        return nil, err
    }
    return &model.Category{ID: c.ID, UserID: c.UserId, Name: c.Name, Type: c.Type, Color: c.Color, Icon: c.Icon}, nil
}

func (r *repoImpl) DeleteCategory(id string) error {
    _, err := r.prisma.Prisma.Category.DeleteOne(
        prisma.Category.ID.Equals(id),
    ).Exec(context.Background())
    return err
}
