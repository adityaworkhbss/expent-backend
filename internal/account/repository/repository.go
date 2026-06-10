package account

import (
    "context"
    "expent-backend/internal/account/model"
    "expent-backend/internal/infrastructure/prisma"
    "github.com/steebchen/prisma-client-go/runtime"
)

type Repository interface {
    ListAccounts(userID string) ([]model.Account, error)
    CreateAccount(acc model.Account) (*model.Account, error)
    GetAccountByID(id string) (*model.Account, error)
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
        prisma.Account.UserId.Equals(userID),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    var result []model.Account
    for _, a := range accs {
        result = append(result, model.Account{ID: a.ID, UserID: a.UserId, Name: a.Name, Type: a.Type, Balance: a.Balance})
    }
    return result, nil
}

func (r *repoImpl) CreateAccount(acc model.Account) (*model.Account, error) {
    a, err := r.prisma.Prisma.Account.CreateOne(
        prisma.Account.UserId.Set(acc.UserID),
        prisma.Account.Name.Set(acc.Name),
        prisma.Account.Type.Set(acc.Type),
        prisma.Account.Balance.Set(acc.Balance),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    return &model.Account{ID: a.ID, UserID: a.UserId, Name: a.Name, Type: a.Type, Balance: a.Balance}, nil
}

func (r *repoImpl) GetAccountByID(id string) (*model.Account, error) {
    a, err := r.prisma.Prisma.Account.FindFirst(
        prisma.Account.ID.Equals(id),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        return nil, err
    }
    return &model.Account{ID: a.ID, UserID: a.UserId, Name: a.Name, Type: a.Type, Balance: a.Balance}, nil
}

func (r *repoImpl) DeleteAccount(id string) error {
    _, err := r.prisma.Prisma.Account.DeleteOne(
        prisma.Account.ID.Equals(id),
    ).Exec(context.Background())
    return err
}
