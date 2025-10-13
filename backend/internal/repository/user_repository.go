package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// Count returns number of users in the repository
	Count(ctx context.Context) (int64, error)
	// List returns all users (or a subset) for autocomplete/lookup
	List(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, u *models.User) error
}
