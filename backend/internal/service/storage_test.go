package service_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"example.com/defect-control-system/internal/handler"
	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/service"
)

func TestLocalStorage_SaveFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "uploads_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	viper.Set("uploads.path", tmpDir)
	viper.Set("uploads.max_size", 1024*1024)
	viper.Set("uploads.allowed_types", []string{"image/jpeg", "image/png"})

	// create multipart body
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("files", "test.jpg")
	assert.NoError(t, err)
	_, err = io.Copy(fw, bytes.NewReader([]byte("\xff\xd8\xff\xdb")))
	assert.NoError(t, err)
	w.Close()

	// parse multipart to get FileHeader
	req := httptest.NewRequest(http.MethodPost, "/", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	mr, err := req.MultipartReader()
	assert.NoError(t, err)
	part, err := mr.NextPart()
	assert.NoError(t, err)
	// we need a FileHeader; create temp file and create a FileHeader by reading part into buffer
	tmpf := filepath.Join(tmpDir, "tmp.jpg")
	data, _ := io.ReadAll(part)
	os.WriteFile(tmpf, data, 0644)
	// construct a fake FileHeader by creating a multipart form and parsing it back
	var b2 bytes.Buffer
	w2 := multipart.NewWriter(&b2)
	fw2, _ := w2.CreateFormFile("files", "test.jpg")
	fw2.Write(data)
	w2.Close()
	req2 := httptest.NewRequest(http.MethodPost, "/", &b2)
	req2.Header.Set("Content-Type", w2.FormDataContentType())
	// parse multipart form to populate MultipartForm field
	err = req2.ParseMultipartForm(10 << 20)
	assert.NoError(t, err)
	form := req2.MultipartForm
	fh := form.File["files"][0]

	st := service.NewLocalStorage()
	rel, size, err := st.SaveFile(fh)
	assert.NoError(t, err)
	assert.True(t, size > 0)
	// check file exists
	full := filepath.Join(tmpDir, rel)
	_, err = os.Stat(full)
	assert.NoError(t, err)
}

func TestDownloadHandler(t *testing.T) {
	// prepare temp dir and file
	tmpDir, err := os.MkdirTemp("", "uploads_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	viper.Set("uploads.path", tmpDir)

	// create file
	dir := filepath.Join(tmpDir, "2025", "10", "11")
	os.MkdirAll(dir, 0755)
	fname := filepath.Join(dir, "file.jpg")
	content := []byte("hello")
	os.WriteFile(fname, content, 0644)

	// mock repo and services
	ar := &mockAttachRepoFile{path: filepath.Join("2025", "10", "11", "file.jpg"), fname: "file.jpg"}
	ds := &mockDefectSvc{}
	storage := service.NewLocalStorage()
	h := handler.NewAttachmentHandler(storage, ar, ds)

	r := gin.Default()
	// test-only middleware to inject authenticated user context
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "engineer")
		c.Next()
	})
	r.GET("/attachments/:id", h.Download)
	req := httptest.NewRequest(http.MethodGet, "/attachments/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(content), rr.Body.String())
}

// mock repo returning our file
type mockAttachRepoFile struct{ path, fname string }

func (m *mockAttachRepoFile) Create(ctx context.Context, a *models.Attachment) error { return nil }
func (m *mockAttachRepoFile) FindByID(ctx context.Context, id uint) (*models.Attachment, error) {
	return &models.Attachment{ID: id, Path: m.path, Filename: m.fname, ContentType: "text/plain", UploaderID: 1}, nil
}
func (m *mockAttachRepoFile) ListByDefect(ctx context.Context, defectID uint) ([]*models.Attachment, error) {
	return []*models.Attachment{{ID: 1, Path: m.path, Filename: m.fname}}, nil
}

// reuse mockDefectSvc from handler tests
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
	return &models.Defect{ID: id, Title: "mock"}, nil
}
