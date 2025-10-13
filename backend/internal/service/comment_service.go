package service

import (
	"context"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
)

type CreateCommentDTO struct {
	DefectID uint   `json:"defect_id" validate:"required"`
	Body     string `json:"body" validate:"required"`
}

type CommentService interface {
	Create(ctx context.Context, authorID uint, dto CreateCommentDTO) (*models.Comment, error)
	ListByDefect(ctx context.Context, defectID uint) ([]*models.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) CommentService {
	return &commentService{repo: r}
}

func (s *commentService) Create(ctx context.Context, authorID uint, dto CreateCommentDTO) (*models.Comment, error) {
	var authorPtr *uint
	if authorID != 0 {
		v := authorID
		authorPtr = &v
	}
	c := &models.Comment{DefectID: dto.DefectID, AuthorID: authorPtr, Body: dto.Body}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	// reload with author preloaded
	saved, err := s.repo.FindByID(ctx, c.ID)
	if err != nil {
		return c, nil // return created comment even if reload fails
	}
	return saved, nil
}

func (s *commentService) ListByDefect(ctx context.Context, defectID uint) ([]*models.Comment, error) {
	return s.repo.ListByDefect(ctx, defectID)
}
