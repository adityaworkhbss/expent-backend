package budget

import (
    "context"
    "expent-backend/internal/budget/model"
    "expent-backend/internal/infrastructure/prisma"
    "github.com/steebchen/prisma-client-go/runtime"
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
        prisma.Budget.UserId.Equals(userID),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    var result []model.Budget
    for _, b := range bs {
        result = append(result, model.Budget{ID: b.ID, UserID: b.UserId, CategoryID: b.CategoryId, Amount: b.Amount, Period: b.Period})
    }
    return result, nil
}

func (r *repoImpl) CreateBudget(budget model.Budget) (*model.Budget, error) {
    b, err := r.prisma.Prisma.Budget.CreateOne(
        prisma.Budget.UserId.Set(budget.UserID),
        prisma.Budget.CategoryId.Set(budget.CategoryID),
        prisma.Budget.Amount.Set(budget.Amount),
        prisma.Budget.Period.Set(budget.Period),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    return &model.Budget{ID: b.ID, UserID: b.UserId, CategoryID: b.CategoryId, Amount: b.Amount, Period: b.Period}, nil
}

func (r *repoImpl) GetBudgetByID(id string) (*model.Budget, error) {
    b, err := r.prisma.Prisma.Budget.FindFirst(
        prisma.Budget.ID.Equals(id),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        return nil, err
    }
    return &model.Budget{ID: b.ID, UserID: b.UserId, CategoryID: b.CategoryId, Amount: b.Amount, Period: b.Period}, nil
}

func (r *repoImpl) DeleteBudget(id string) error {
    _, err := r.prisma.Prisma.Budget.DeleteOne(
        prisma.Budget.ID.Equals(id),
    ).Exec(context.Background())
    return err
}
