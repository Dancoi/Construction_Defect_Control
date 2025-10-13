package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"gorm.io/gorm"
)

type attachmentRepoPG struct{ db *gorm.DB }

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository { return &attachmentRepoPG{db: db} }

func (r *attachmentRepoPG) Create(ctx context.Context, a *models.Attachment) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *attachmentRepoPG) FindByID(ctx context.Context, id uint) (*models.Attachment, error) {
	var a models.Attachment
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *attachmentRepoPG) ListByDefect(ctx context.Context, defectID uint) ([]*models.Attachment, error) {
	var list []*models.Attachment
	if err := r.db.WithContext(ctx).Where("defect_id = ?", defectID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
