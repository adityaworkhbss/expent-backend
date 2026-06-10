package auth

import (
    "context"
    "expent-backend/internal/auth/model"
    "expent-backend/internal/infrastructure/prisma"
    "github.com/steebchen/prisma-client-go/runtime"
    "go.uber.org/zap"
)

// Repository defines data access methods for auth module.
type Repository interface {
    GetUserByEmail(email string) (*model.User, error)
    CreateUser(email, name string) (*model.User, error)
    StoreRefreshToken(userID, token string) error
    GetRefreshToken(token string) (*model.RefreshToken, error)
}

type repoImpl struct {
    prisma *prisma.PrismaClient
}

func NewRepository(prismaClient *prisma.PrismaClient) Repository {
    return &repoImpl{prisma: prismaClient}
}

func (r *repoImpl) GetUserByEmail(email string) (*model.User, error) {
    user, err := r.prisma.Prisma.User.FindFirst(
        prisma.User.Email.Equals(email),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        zap.S().Error("failed to fetch user", zap.Error(err))
        return nil, err
    }
    return &model.User{ID: user.ID, Email: user.Email, Name: user.Name}, nil
}

func (r *repoImpl) CreateUser(email, name string) (*model.User, error) {
    usr, err := r.prisma.Prisma.User.CreateOne(
        prisma.User.Email.Set(email),
        prisma.User.Name.Set(name),
    ).Exec(context.Background())
    if err != nil {
        zap.S().Error("failed to create user", zap.Error(err))
        return nil, err
    }
    return &model.User{ID: usr.ID, Email: usr.Email, Name: usr.Name}, nil
}

func (r *repoImpl) StoreRefreshToken(userID, token string) error {
    _, err := r.prisma.Prisma.RefreshToken.CreateOne(
        prisma.RefreshToken.Token.Set(token),
        prisma.RefreshToken.ExpiresAt.Set(runtime.Now().Add(runtime.Duration(7*24*60*60*1000))), // placeholder 7 days
        prisma.RefreshToken.UserId.Set(userID),
    ).Exec(context.Background())
    if err != nil {
        zap.S().Error("failed to store refresh token", zap.Error(err))
    }
    return err
}

func (r *repoImpl) GetRefreshToken(token string) (*model.RefreshToken, error) {
    rt, err := r.prisma.Prisma.RefreshToken.FindFirst(
        prisma.RefreshToken.Token.Equals(token),
    ).Exec(context.Background())
    if err != nil {
        if runtime.IsNotFound(err) {
            return nil, nil
        }
        zap.S().Error("failed to fetch refresh token", zap.Error(err))
        return nil, err
    }
    return &model.RefreshToken{ID: rt.ID, Token: rt.Token, ExpiresAt: rt.ExpiresAt, UserID: rt.UserId}, nil
}
