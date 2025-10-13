package repository

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"gorm.io/gorm"
)

type userRepoPG struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepoPG{db: db}
}

func (r *userRepoPG) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepoPG) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepoPG) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepoPG) Count(ctx context.Context) (int64, error) {
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *userRepoPG) List(ctx context.Context) ([]*models.User, error) {
	var list []*models.User
	// select a small set of fields for performance
	if err := r.db.WithContext(ctx).Select("id", "name", "email", "role").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepoPG) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}
