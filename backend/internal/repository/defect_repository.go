package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
)

type DefectRepository interface {
	Create(ctx context.Context, d *models.Defect) error
	FindByID(ctx context.Context, id uint) (*models.Defect, error)
	ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error)
	Update(ctx context.Context, d *models.Defect) error
}
