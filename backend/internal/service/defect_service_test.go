package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/service"
)

type mockDefectRepo struct{}

func (m *mockDefectRepo) Create(ctx context.Context, d *models.Defect) error { d.ID = 1; return nil }
func (m *mockDefectRepo) FindByID(ctx context.Context, id uint) (*models.Defect, error) {
	return &models.Defect{ID: id}, nil
}
func (m *mockDefectRepo) ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error) {
	return []*models.Defect{}, nil
}

type mockProjectRepoNotFound struct{}

func (m *mockProjectRepoNotFound) Create(ctx context.Context, p *models.Project) error {
	return errors.New("not implemented")
}
func (m *mockProjectRepoNotFound) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	return nil, errors.New("not found")
}
func (m *mockProjectRepoNotFound) List(ctx context.Context) ([]*models.Project, error) {
	return nil, errors.New("not implemented")
}

func TestCreateDefect_ProjectNotFound(t *testing.T) {
	repo := &mockDefectRepo{}
	projRepo := &mockProjectRepoNotFound{}
	s := service.NewDefectService(repo, projRepo)
	dto := service.CreateDefectDTO{ProjectID: 999, Title: "x"}
	_, err := s.Create(context.Background(), dto)
	assert.Error(t, err)
}

func TestCreateDefect_Success(t *testing.T) {
	repo := &mockDefectRepo{}
	projRepo := &mockProjectRepo{}
	s := service.NewDefectService(repo, projRepo)
	dto := service.CreateDefectDTO{ProjectID: 1, Title: "x"}
	d, err := s.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), d.ID)
}

// mockProjectRepo to satisfy successful path
type mockProjectRepo struct{}

func (m *mockProjectRepo) Create(ctx context.Context, p *models.Project) error { p.ID = 1; return nil }
func (m *mockProjectRepo) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	return &models.Project{ID: id}, nil
}
func (m *mockProjectRepo) List(ctx context.Context) ([]*models.Project, error) {
	return []*models.Project{}, nil
}
