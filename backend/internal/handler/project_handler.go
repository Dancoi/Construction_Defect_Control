package handler

import (
	"fmt"
	"net/http"

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
	var dto service.CreateDefectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
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
	dto.ProjectID = projectID
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

// ListDefects godoc
// @Summary List defects for a project
// @Description Get defects for given project id
// @Tags defects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} handler.DefectResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/defects [get]
