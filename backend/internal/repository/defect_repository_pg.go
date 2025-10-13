package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"gorm.io/gorm"
)

type defectRepoPG struct{ db *gorm.DB }

func NewDefectRepository(db *gorm.DB) DefectRepository { return &defectRepoPG{db: db} }

func (r *defectRepoPG) Create(ctx context.Context, d *models.Defect) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *defectRepoPG) FindByID(ctx context.Context, id uint) (*models.Defect, error) {
	var d models.Defect
	if err := r.db.WithContext(ctx).First(&d, id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *defectRepoPG) ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error) {
	var list []*models.Defect
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *defectRepoPG) Update(ctx context.Context, d *models.Defect) error {
	return r.db.WithContext(ctx).Save(d).Error
}
