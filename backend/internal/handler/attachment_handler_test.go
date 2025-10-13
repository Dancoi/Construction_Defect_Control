package handler_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	hpkg "example.com/defect-control-system/internal/handler"
	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/service"
)

// mock attachment repo that records created attachments
type mockAttachRepo struct{}

func (m *mockAttachRepo) Create(ctx context.Context, a *models.Attachment) error {
	a.ID = 1
	return nil
}
func (m *mockAttachRepo) FindByID(ctx context.Context, id uint) (*models.Attachment, error) {
	return &models.Attachment{ID: id, Path: "2025/10/11/file.jpg", Filename: "file.jpg", ContentType: "image/jpeg"}, nil
}
func (m *mockAttachRepo) ListByDefect(ctx context.Context, defectID uint) ([]*models.Attachment, error) {
	return []*models.Attachment{}, nil
}

// mock defect service that accepts any project id
type mockDefectSvc struct{}

func (m *mockDefectSvc) Create(ctx context.Context, dto service.CreateDefectDTO) (*models.Defect, error) {
	return &models.Defect{ID: 1, Title: dto.Title}, nil
}
func (m *mockDefectSvc) ListByProject(ctx context.Context, projectID uint) ([]*models.Defect, error) {
	return []*models.Defect{}, nil
}
func (m *mockDefectSvc) FindByID(ctx context.Context, id uint) (*models.Defect, error) {
	return &models.Defect{ID: id, Title: "mock"}, nil
}
func (m *mockDefectSvc) Update(ctx context.Context, id uint, dto service.UpdateDefectDTO) (*models.Defect, error) {
	// for tests, just return a defect with updated title if provided
	d := &models.Defect{ID: id}
	if dto.Title != nil {
		d.Title = *dto.Title
	} else {
		d.Title = "mock"
	}
	return d, nil
}

func TestUploadHandler(t *testing.T) {
	// temp uploads dir
	tmpDir, err := os.MkdirTemp("", "uploads_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	viper.Set("uploads.path", tmpDir)
	viper.Set("uploads.allowed_types", []string{"image/jpeg", "image/png"})

	storage := service.NewLocalStorage()
	attachRepo := &mockAttachRepo{}
	defectSvc := &mockDefectSvc{}
	h := hpkg.NewAttachmentHandler(storage, attachRepo, defectSvc)

	r := gin.Default()
	r.POST("/upload/:id", h.Upload)

	// create multipart body
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("files", "test.jpg")
	assert.NoError(t, err)
	_, err = io.Copy(fw, bytes.NewReader([]byte("\xff\xd8\xff\xdb")))
	assert.NoError(t, err)
	w.Close()

	req := httptest.NewRequest("POST", "/upload/1", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}
