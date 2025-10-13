package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/defect-control-system/internal/repository"
	"example.com/defect-control-system/internal/service"
)

type UserHandler struct {
	repo repository.UserRepository
	auth service.AuthService
}

func NewUserHandler(r repository.UserRepository, a service.AuthService) *UserHandler {
	return &UserHandler{repo: r, auth: a}
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

// Me godoc
// @Summary Current user
// @Description Get current authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} object
// @Security BearerAuth
// @Router /api/v1/users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "unauthenticated"})
		return
	}
	id, ok := uid.(uint)
	if !ok {
		if if64, ok2 := uid.(int64); ok2 {
			id = uint(if64)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "invalid user id"})
			return
		}
	}
	u, err := h.auth.Me(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": u})
}

// UpdateMe godoc
// @Summary Update current user
// @Description Update current user's name/email
// @Tags users
// @Accept json
// @Produce json
// @Param body body map[string]string true "Update"
// @Success 200 {object} object
// @Security BearerAuth
// @Router /api/v1/users/me [patch]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "unauthenticated"})
		return
	}
	id, ok := uid.(uint)
	if !ok {
		if if64, ok2 := uid.(int64); ok2 {
			id = uint(if64)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "invalid user id"})
			return
		}
	}
	var body struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	u, err := h.auth.UpdateProfile(c.Request.Context(), id, body.Name, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": u})
}

// UpdateUser godoc
// @Summary Update user (admin)
// @Description Update a user's role/name/email (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body map[string]string true "Update"
// @Success 200 {object} object
// @Security BearerAuth
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	uid64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid id"})
		return
	}
	id := uint(uid64)

	var body struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
		Role  *string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	u, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "user not found"})
		return
	}
	if body.Name != nil {
		u.Name = *body.Name
	}
	if body.Email != nil {
		u.Email = *body.Email
	}
	if body.Role != nil {
		u.Role = *body.Role
	}
	if err := h.repo.Update(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role}})
}
