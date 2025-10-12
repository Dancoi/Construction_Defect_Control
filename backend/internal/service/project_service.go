package service

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
)

type CreateProjectDTO struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address"`
}

type ProjectService interface {
	Create(ctx context.Context, dto CreateProjectDTO) (*models.Project, error)
	GetByID(ctx context.Context, id uint) (*models.Project, error)
	List(ctx context.Context) ([]*models.Project, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(r repository.ProjectRepository) ProjectService {
	return &projectService{repo: r}
}

func (s *projectService) Create(ctx context.Context, dto CreateProjectDTO) (*models.Project, error) {
	p := &models.Project{
		Name:    dto.Name,
		Address: dto.Address,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *projectService) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *projectService) List(ctx context.Context) ([]*models.Project, error) {
	return s.repo.List(ctx)
}
