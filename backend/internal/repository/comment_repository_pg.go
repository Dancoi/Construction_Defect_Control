package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"gorm.io/gorm"
)

type commentRepoPG struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) CommentRepository { return &commentRepoPG{db: db} }

func (r *commentRepoPG) Create(ctx context.Context, c *models.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *commentRepoPG) ListByDefect(ctx context.Context, defectID uint) ([]*models.Comment, error) {
	var list []*models.Comment
	if err := r.db.WithContext(ctx).Where("defect_id = ?", defectID).Preload("Author").Order("created_at asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *commentRepoPG) FindByID(ctx context.Context, id uint) (*models.Comment, error) {
	var c models.Comment
	if err := r.db.WithContext(ctx).Preload("Author").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
