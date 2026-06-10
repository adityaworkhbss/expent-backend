package transaction

import (
    "context"
    "expent-backend/internal/transaction/model"
    "expent-backend/internal/infrastructure/prisma"
    "github.com/steebchen/prisma-client-go/runtime"
)

type Repository interface {
    ListTransactions(userID string) ([]model.Transaction, error)
    CreateTransaction(tx model.Transaction) (*model.Transaction, error)
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
    txs, err := r.prisma.Prisma.Transaction.FindMany(
        prisma.Transaction.UserId.Equals(userID),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    var result []model.Transaction
    for _, t := range txs {
        result = append(result, model.Transaction{ID: t.ID, UserID: t.UserId, AccountID: t.AccountId, CategoryID: t.CategoryId, Amount: t.Amount, Type: t.Type, Timestamp: t.Timestamp, Description: t.Description})
    }
    return result, nil
}

func (r *repoImpl) CreateTransaction(tx model.Transaction) (*model.Transaction, error) {
    t, err := r.prisma.Prisma.Transaction.CreateOne(
        prisma.Transaction.UserId.Set(tx.UserID),
        prisma.Transaction.AccountId.Set(tx.AccountID),
        prisma.Transaction.CategoryId.Set(tx.CategoryID),
        prisma.Transaction.Amount.Set(tx.Amount),
        prisma.Transaction.Type.Set(tx.Type),
        prisma.Transaction.Timestamp.Set(tx.Timestamp),
        prisma.Transaction.Description.Set(tx.Description),
    ).Exec(context.Background())
    if err != nil {
        return nil, err
    }
    return &model.Transaction{ID: t.ID, UserID: t.UserId, AccountID: t.AccountId, CategoryID: t.CategoryId, Amount: t.Amount, Type: t.Type, Timestamp: t.Timestamp, Description: t.Description}, nil
}

func (r *repoImpl) GetTransactionByID(id string) (*model.Transaction, error) {
    t, err := r.prisma.Prisma.Transaction.FindFirst(
        prisma.Transaction.ID.Equals(id),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        return nil, err
    }
    return &model.Transaction{ID: t.ID, UserID: t.UserId, AccountID: t.AccountId, CategoryID: t.CategoryId, Amount: t.Amount, Type: t.Type, Timestamp: t.Timestamp, Description: t.Description}, nil
}

func (r *repoImpl) DeleteTransaction(id string) error {
    _, err := r.prisma.Prisma.Transaction.DeleteOne(
        prisma.Transaction.ID.Equals(id),
    ).Exec(context.Background())
    return err
}
