package service

import (
	"context"
	"errors"
	"time"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
)

type CreateDefectDTO struct {
	ProjectID   uint       `json:"project_id" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"`
	AssigneeID  uint       `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    string     `json:"priority"`
}

type DefectService interface {
	Create(ctx context.Context, dto CreateDefectDTO) (*models.Defect, error)
	ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error)
	FindByID(ctx context.Context, id uint) (*models.Defect, error)
	Update(ctx context.Context, id uint, dto UpdateDefectDTO) (*models.Defect, error)
}

type defectService struct {
	repo        repository.DefectRepository
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

func NewDefectService(r repository.DefectRepository, pr repository.ProjectRepository, ur repository.UserRepository) DefectService {
	return &defectService{repo: r, projectRepo: pr, userRepo: ur}
}

func (s *defectService) Create(ctx context.Context, dto CreateDefectDTO) (*models.Defect, error) {
	// ensure project exists
	if _, err := s.projectRepo.FindByID(ctx, dto.ProjectID); err != nil {
		return nil, errors.New("project not found")
	}

	var assigneePtr *uint
	if dto.AssigneeID != 0 {
		// validate assignee exists
		if _, err := s.userRepo.FindByID(ctx, dto.AssigneeID); err != nil {
			return nil, errors.New("assignee not found")
		}
		v := dto.AssigneeID
		assigneePtr = &v
	}
	d := &models.Defect{
		ProjectID:   dto.ProjectID,
		Title:       dto.Title,
		Description: dto.Description,
		Severity:    dto.Severity,
		Status:      "open",
		AssigneeID:  assigneePtr,
		DueDate:     dto.DueDate,
		Priority:    dto.Priority,
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

type UpdateDefectDTO struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Severity    *string    `json:"severity"`
	AssigneeID  *uint      `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Priority    *string    `json:"priority"`
	Status      *string    `json:"status"`
}

func (s *defectService) Update(ctx context.Context, id uint, dto UpdateDefectDTO) (*models.Defect, error) {
	d, err := s.repo.FindByID(ctx, id)
	if err != nil || d == nil {
		return nil, errors.New("defect not found")
	}
	if dto.Title != nil {
		d.Title = *dto.Title
	}
	if dto.Description != nil {
		d.Description = *dto.Description
	}
	if dto.Severity != nil {
		d.Severity = *dto.Severity
	}
	if dto.AssigneeID != nil {
		if *dto.AssigneeID == 0 {
			d.AssigneeID = nil
		} else {
			if _, err := s.userRepo.FindByID(ctx, *dto.AssigneeID); err != nil {
				return nil, errors.New("assignee not found")
			}
			v := *dto.AssigneeID
			d.AssigneeID = &v
		}
	}
	if dto.DueDate != nil {
		d.DueDate = dto.DueDate
	}
	if dto.Priority != nil {
		d.Priority = *dto.Priority
	}
	if dto.Status != nil {
		d.Status = *dto.Status
	}
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *defectService) ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *defectService) FindByID(ctx context.Context, id uint) (*models.Defect, error) {
	return s.repo.FindByID(ctx, id)
}
