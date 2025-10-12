package service

import (
	"context"
	"errors"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
)

type CreateDefectDTO struct {
	ProjectID   uint   `json:"project_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type DefectService interface {
	Create(ctx context.Context, dto CreateDefectDTO) (*models.Defect, error)
	ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error)
}

type defectService struct {
	repo        repository.DefectRepository
	projectRepo repository.ProjectRepository
}

func NewDefectService(r repository.DefectRepository, pr repository.ProjectRepository) DefectService {
	return &defectService{repo: r, projectRepo: pr}
}

func (s *defectService) Create(ctx context.Context, dto CreateDefectDTO) (*models.Defect, error) {
	// ensure project exists
	if _, err := s.projectRepo.FindByID(ctx, dto.ProjectID); err != nil {
		return nil, errors.New("project not found")
	}

	d := &models.Defect{
		ProjectID:   dto.ProjectID,
		Title:       dto.Title,
		Description: dto.Description,
		Severity:    dto.Severity,
		Status:      "open",
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *defectService) ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error) {
	return s.repo.ListByProject(ctx, projectID)
}
