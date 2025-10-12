package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/service"
)

type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, u *models.User) error { u.ID = 1; return nil }
func (m *mockUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return &models.User{ID: id, Email: "a@b.com", PasswordHash: "", Name: "x"}, nil
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Count(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserRepo) List(ctx context.Context) ([]*models.User, error) {
	return []*models.User{{ID: 1, Name: "Test", Email: "t@example.com"}}, nil
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepo{}
	s := service.NewAuthService(repo)
	dto := service.RegisterDTO{Name: "Test", Email: "t@example.com", Password: "password123"}
	u, err := s.Register(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), u.ID)
}
