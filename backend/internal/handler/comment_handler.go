package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/defect-control-system/internal/service"
)

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(s service.CommentService) *CommentHandler { return &CommentHandler{svc: s} }

// CreateComment godoc
// @Summary Create a comment for a defect
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param defectId path int true "Defect ID"
// @Param body body service.CreateCommentDTO true "Create Comment"
// @Success 201 {object} handler.CommentResponse
// @Security BearerAuth
// @Router /api/v1/projects/{id}/defects/{defectId}/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	// parse defect id from path to ensure it is present
	did := c.Param("defectId")
	var defectID uint
	if _, err := fmt.Sscanf(did, "%d", &defectID); err != nil || defectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id"})
		return
	}
	var dto service.CreateCommentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	// set defect id from path
	dto.DefectID = defectID
	authorID := uint(0)
	if v, ok := c.Get("user_id"); ok {
		if uid, ok2 := v.(uint); ok2 {
			authorID = uid
		}
	}
	cm, err := h.svc.Create(c.Request.Context(), authorID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "data": cm})
}

// ListComments godoc
// @Summary List comments for a defect
// @Produce json
// @Param id path int true "Project ID"
// @Param defectId path int true "Defect ID"
// @Success 200 {array} handler.CommentResponse
// @Router /api/v1/projects/{id}/defects/{defectId}/comments [get]
func (h *CommentHandler) List(c *gin.Context) {
	// prefer path param defectId, but allow query param defect_id for flexibility
	did := c.Param("defectId")
	var defectID uint
	if did != "" {
		if _, err := fmt.Sscanf(did, "%d", &defectID); err != nil || defectID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id"})
			return
		}
	} else {
		// fallback to query param
		q := c.Query("defect_id")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "defect_id required"})
			return
		}
		if _, err := fmt.Sscanf(q, "%d", &defectID); err != nil || defectID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid defect id"})
			return
		}
	}
	list, err := h.svc.ListByDefect(c.Request.Context(), defectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	// simple logging for debugging: show how many comments returned
	fmt.Printf("ListComments: defect=%d count=%d\n", defectID, len(list))
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": list})
}
