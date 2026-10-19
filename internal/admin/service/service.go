package service

import (
	"context"
	"errors"

	"expent-backend/internal/admin/repository"
	"expent-backend/internal/infrastructure/jwt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAdminAlreadyExists = errors.New("an admin account already exists; setup is a one-time operation")

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, email, password string) (accessToken string, admin *repository.AdminResponse, err error) {
	adminModel, err := s.repo.GetAdminByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if adminModel == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminModel.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	accessToken, err = jwt.GenerateAdminAccessToken(adminModel.ID)
	if err != nil {
		return "", nil, err
	}

	return accessToken, &repository.AdminResponse{
		ID:    adminModel.ID,
		Email: adminModel.Email,
		Name:  adminModel.Name,
	}, nil
}

// SetupAdmin creates the very first admin account. It is a one-time
// operation: once any admin exists, this always fails, so the endpoint
// cannot be used to mint additional accounts later.
func (s *Service) SetupAdmin(ctx context.Context, email, password, name string) (*repository.AdminResponse, error) {
	count, err := s.repo.CountAdmins()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrAdminAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin, err := s.repo.CreateAdmin(email, string(hash), name)
	if err != nil {
		return nil, err
	}

	return &repository.AdminResponse{
		ID:    admin.ID,
		Email: admin.Email,
		Name:  admin.Name,
	}, nil
}

func (s *Service) ValidateAdminToken(tokenString string) (string, error) {
	claims, err := jwt.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
