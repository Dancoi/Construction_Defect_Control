package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"gorm.io/gorm"
)

type projectRepoPG struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) ProjectRepository { return &projectRepoPG{db: db} }

func (r *projectRepoPG) Create(ctx context.Context, p *models.Project) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *projectRepoPG) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	var p models.Project
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepoPG) List(ctx context.Context) ([]*models.Project, error) {
	var list []*models.Project
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *projectRepoPG) Update(ctx context.Context, p *models.Project) error {
	return r.db.WithContext(ctx).Save(p).Error
}
