package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"example.com/defect-control-system/internal/models"
	"example.com/defect-control-system/internal/repository"
	"example.com/defect-control-system/internal/service"
)

type AttachmentHandler struct {
	storage    service.StorageService
	attachRepo repository.AttachmentRepository
	defectSvc  service.DefectService
}

func NewAttachmentHandler(s service.StorageService, ar repository.AttachmentRepository, ds service.DefectService) *AttachmentHandler {
	return &AttachmentHandler{storage: s, attachRepo: ar, defectSvc: ds}
}

// UploadAttachments godoc
// @Summary Upload attachments to defect
// @Description Upload one or multiple files as attachments to a defect
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Defect ID"
// @Param files formData file true "files"
// @Success 201 {array} handler.AttachmentResponse
// @Security BearerAuth
// @Router /api/v1/projects/{id}/attachments [post]
func (h *AttachmentHandler) Upload(c *gin.Context) {
	// The route is mounted under /projects/:id/attachments but historically the handler
	// treated the path param as a defect id. Frontend sends defect id as query param
	// (defect_id). To be robust we accept defect_id from query or form field. If
	// absent, fall back to using the path param as defect id.
	idParam := c.Param("id")
	var pathID uint
	if _, err := fmt.Sscanf(idParam, "%d", &pathID); err != nil {
		pathID = 0
	}

	// prefer defect_id query param
	defectIDStr := c.Query("defect_id")
	var defectID uint
	if defectIDStr != "" {
		if _, err := fmt.Sscanf(defectIDStr, "%d", &defectID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect_id query param"})
			return
		}
	} else {
		// also accept form field "defect_id" when multipart
		if v := c.PostForm("defect_id"); v != "" {
			if _, err := fmt.Sscanf(v, "%d", &defectID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect_id form field"})
				return
			}
		}
	}

	// if still no defect id, fall back to using path param (backward compatibility)
	if defectID == 0 && pathID != 0 {
		defectID = pathID
	}

	if defectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "defect_id is required"})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	files := form.File["files"]
	var results []gin.H
	for _, fh := range files {
		relpath, size, err := h.storage.SaveFile(fh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		uploaderID := uint(0)
		if v, ok := c.Get("user_id"); ok {
			if uid, ok2 := v.(uint); ok2 {
				uploaderID = uid
			}
		}
		a := &models.Attachment{
			DefectID:    defectID,
			UploaderID:  uploaderID,
			Path:        relpath,
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        size,
		}
		if err := h.attachRepo.Create(c.Request.Context(), a); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		results = append(results, gin.H{"id": a.ID, "filename": a.Filename, "url": filepath.Join("/uploads", a.Path), "content_type": a.ContentType, "size": a.Size})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "data": results})
}

// DownloadAttachment godoc
// @Summary Download attachment
// @Description Download attachment by id
// @Tags attachments
// @Produce application/octet-stream
// @Param id path int true "Attachment ID"
// @Success 200 {object} handler.AttachmentResponse
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/attachments/{id} [get]
func (h *AttachmentHandler) Download(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid id"})
		return
	}
	a, err := h.attachRepo.FindByID(c.Request.Context(), id)
	if err != nil || a == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "attachment not found"})
		return
	}
	// check ownership/permission: allow if uploader, or role in manager/admin/stakeholder
	allowed := false
	if v, ok := c.Get("user_id"); ok {
		if uid, ok2 := v.(uint); ok2 && uid == a.UploaderID {
			allowed = true
		}
	}
	if !allowed {
		if rv, ok := c.Get("role"); ok {
			if role, ok2 := rv.(string); ok2 {
				if role == "manager" || role == "admin" || role == "stakeholder" {
					allowed = true
				}
			}
		}
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "error": "forbidden"})
		return
	}
	base := viper.GetString("uploads.path")
	if base == "" {
		base = "./uploads"
	}
	// build full path
	full := filepath.Join(base, a.Path)
	// ensure file exists
	if _, err := os.Stat(full); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "file not found"})
		return
	}
	c.Header("Content-Type", a.ContentType)
	// If content type is an image, prefer inline display in browser
	if a.ContentType != "" && (a.ContentType == "image/jpeg" || a.ContentType == "image/png" || a.ContentType == "image/gif" || a.ContentType == "image/webp") {
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", a.Filename))
	} else {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", a.Filename))
	}
	c.File(full)
}

// ListAttachments godoc
// @Summary List attachments
// @Description List attachments by defect id
// @Tags attachments
// @Produce json
// @Param defect_id query int false "Defect ID"
// @Param id path int false "Project or defect id in path"
// @Success 200 {array} handler.AttachmentResponse
// @Router /api/v1/attachments [get]
func (h *AttachmentHandler) List(c *gin.Context) {
	// prefer query param defect_id
	q := c.Query("defect_id")
	var defectID uint
	if q != "" {
		if _, err := fmt.Sscanf(q, "%d", &defectID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect_id"})
			return
		}
	} else {
		// allow path param id when route is /projects/:id/defects/:defectId/attachments
		did := c.Param("defectId")
		if did != "" {
			if _, err := fmt.Sscanf(did, "%d", &defectID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id in path"})
				return
			}
		}
	}
	if defectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "defect_id required"})
		return
	}
	list, err := h.attachRepo.ListByDefect(c.Request.Context(), defectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	var out []gin.H
	for _, a := range list {
		out = append(out, gin.H{"id": a.ID, "filename": a.Filename, "url": filepath.Join("/uploads", a.Path), "content_type": a.ContentType, "size": a.Size})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": out})
}

// DownloadAttachment godoc
// @Summary Download attachment
// @Description Download attachment by id
// @Tags attachments
// @Produce octet-stream
// @Param id path int true "Attachment ID"
// @Success 200 {object} handler.AttachmentResponse
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /attachments/{id} [get]
