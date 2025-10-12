package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/defect-control-system/internal/repository"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

// ListUsers godoc
// @Summary List users
// @Description Return list of users for lookup/autocomplete
// @Tags users
// @Produce json
// @Success 200 {array} object
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	// return minimal fields
	var out []gin.H
	for _, u := range list {
		out = append(out, gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": out})
}
