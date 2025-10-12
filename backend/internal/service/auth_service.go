package service

import (
	"context"
	"errors"
	"os"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
	"example.com/defect-control-system/internal/utils"
)

// DTOs
type RegisterDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthService defines public methods used by handlers/tests
type AuthService interface {
	Register(ctx context.Context, dto RegisterDTO) (*models.User, error)
	Login(ctx context.Context, dto LoginDTO) (string, *models.User, error)
	Me(ctx context.Context, id uint) (*models.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

// NewAuthService constructs a new AuthService. jwtSecret may be empty for tests; in that case a default is used.
func NewAuthService(r repository.UserRepository) AuthService {
	return &authService{repo: r, jwtSecret: "secret"}
}

// NewAuthServiceWithSecret constructs AuthService with a provided jwt secret.
func NewAuthServiceWithSecret(r repository.UserRepository, secret string) AuthService {
	if secret == "" {
		secret = "secret"
	}
	return &authService{repo: r, jwtSecret: secret}
}

func (s *authService) Register(ctx context.Context, dto RegisterDTO) (*models.User, error) {
	// naive implementation: hash password and create user
	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	// default role
	role := "engineer"

	// bootstrap first admin: if enabled via env/config and no users exist, make first user admin
	if os.Getenv("AUTH_BOOTSTRAP_FIRST_ADMIN") == "true" {
		if cnt, err := s.repo.Count(ctx); err == nil && cnt == 0 {
			role = "admin"
		}
	}

	u := &models.User{
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: hash,
		Role:         role,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *authService) Login(ctx context.Context, dto LoginDTO) (string, *models.User, error) {
	u, err := s.repo.FindByEmail(ctx, dto.Email)
	if err != nil || u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if !utils.ComparePassword(u.PasswordHash, dto.Password) {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := utils.CreateJWT(s.jwtSecret, u)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (s *authService) Me(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}
