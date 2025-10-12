package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
)

type AttachmentRepository interface {
	Create(ctx context.Context, a *models.Attachment) error
	FindByID(ctx context.Context, id uint) (*models.Attachment, error)
}
