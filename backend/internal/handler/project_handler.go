package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"example.com/defect-control-system/internal/service"
)

type ProjectHandler struct {
	svc       service.ProjectService
	defectSvc service.DefectService
}

func NewProjectHandler(s service.ProjectService, d service.DefectService) *ProjectHandler {
	return &ProjectHandler{svc: s, defectSvc: d}
}

// CreateProject godoc
// @Summary Create a project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Param body body service.CreateProjectDTO true "Create Project"
// @Success 201 {object} handler.ProjectResponse
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var dto service.CreateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	p, err := h.svc.Create(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "data": p})
}

// ListProjects godoc
// @Summary List projects
// @Description Get list of projects
// @Tags projects
// @Produce json
// @Success 200 {array} handler.ProjectResponse
// @Router /api/v1/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": list})
}

// CreateDefect godoc
// @Summary Create a defect in a project
// @Description Create a defect under specified project
// @Tags defects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body service.CreateDefectDTO true "Create Defect"
// @Success 201 {object} handler.DefectResponse
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/projects/{id}/defects [post]
func (h *ProjectHandler) CreateDefect(c *gin.Context) {
	// bind into a local request type so we can accept multiple date formats
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Severity    string `json:"severity"`
		AssigneeID  uint   `json:"assignee_id"`
		DueDate     string `json:"due_date"`
		Priority    string `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	// ensure project id comes from path param
	pid := c.Param("id")
	var projectID uint
	if _, err := fmt.Sscanf(pid, "%d", &projectID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid project id"})
		return
	}
	// build service DTO and parse due_date if provided
	var dto service.CreateDefectDTO
	dto.ProjectID = projectID
	dto.Title = req.Title
	dto.Description = req.Description
	dto.Severity = req.Severity
	dto.AssigneeID = req.AssigneeID
	dto.Priority = req.Priority
	if req.DueDate != "" {
		// try several common date formats: RFC3339 and date-only YYYY-MM-DD
		var parsed time.Time
		var parseErr error
		layouts := []string{time.RFC3339, "2006-01-02", "2006-01-02T15:04:05", "2006-01-02 15:04:05"}
		for _, l := range layouts {
			parsed, parseErr = time.Parse(l, req.DueDate)
			if parseErr == nil {
				dto.DueDate = &parsed
				break
			}
		}
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": fmt.Sprintf("invalid due_date format: %v", parseErr)})
			return
		}
	}

	d, err := h.defectSvc.Create(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "data": d})
}

// CreateDefect godoc
// @Summary Create a defect in a project
// @Description Create a defect under specified project
// @Tags defects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body service.CreateDefectDTO true "Create Defect"
// @Success 201 {object} handler.DefectResponse
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /projects/{id}/defects [post]

func (h *ProjectHandler) ListDefects(c *gin.Context) {
	pid := c.Param("id")
	// parse id
	var id uint
	_, err := fmt.Sscanf(pid, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid project id"})
		return
	}
	list, err := h.defectSvc.ListByProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": list})
}

// GetProject godoc
// @Summary Get project
// @Description Get project by id
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} handler.ProjectResponse
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	pid := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(pid, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid project id"})
		return
	}
	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": p})
}

// GetDefect godoc
// @Summary Get defect
// @Description Get defect by id
// @Tags defects
// @Produce json
// @Param id path int true "Project ID"
// @Param defectId path int true "Defect ID"
// @Success 200 {object} handler.DefectResponse
// @Router /api/v1/projects/{id}/defects/{defectId} [get]
func (h *ProjectHandler) GetDefect(c *gin.Context) {
	did := c.Param("defectId")
	var id uint
	if _, err := fmt.Sscanf(did, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id"})
		return
	}
	d, err := h.defectSvc.FindByID(c.Request.Context(), id)
	if err != nil || d == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "defect not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": d})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update project fields
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param body body service.UpdateProjectDTO true "Update Project"
// @Success 200 {object} handler.ProjectResponse
// @Failure 400 {object} map[string]interface{}}
// @Security BearerAuth
// @Router /api/v1/projects/{id} [patch]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	pid := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(pid, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid project id"})
		return
	}
	var dto service.UpdateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": p})
}

// UpdateDefect godoc
// @Summary Update a defect
// @Description Update defect fields
// @Tags defects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param defectId path int true "Defect ID"
// @Param body body service.UpdateDefectDTO true "Update Defect"
// @Success 200 {object} handler.DefectResponse
// @Failure 400 {object} map[string]interface{}}
// @Security BearerAuth
// @Router /api/v1/projects/{id}/defects/{defectId} [patch]
func (h *ProjectHandler) UpdateDefect(c *gin.Context) {
	did := c.Param("defectId")
	var id uint
	if _, err := fmt.Sscanf(did, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id"})
		return
	}
	var dto service.UpdateDefectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	d, err := h.defectSvc.Update(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": d})
}

// ListDefects godoc
// @Summary List defects for a project
// @Description Get defects for given project id
// @Tags defects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} handler.DefectResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/defects [get]
