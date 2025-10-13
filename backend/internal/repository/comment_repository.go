package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
)

type CommentRepository interface {
	Create(ctx context.Context, c *models.Comment) error
	ListByDefect(ctx context.Context, defectID uint) ([]*models.Comment, error)
	FindByID(ctx context.Context, id uint) (*models.Comment, error)
}
