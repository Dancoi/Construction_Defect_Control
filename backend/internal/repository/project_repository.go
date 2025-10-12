package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, p *models.Project) error
	FindByID(ctx context.Context, id uint) (*models.Project, error)
	List(ctx context.Context) ([]*models.Project, error)
}
