package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"example.com/defect-control-system/internal/service"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

// Register godoc
// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.RegisterDTO true "Register"
// @Success 201 {object} handler.AuthResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto service.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if err := validator.New().Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	u, err := h.svc.Register(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "data": u})
}

// Login godoc
// @Summary Login
// @Description Login user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.LoginDTO true "Login"
// @Success 200 {object} handler.AuthResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto service.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	token, user, err := h.svc.Login(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": gin.H{"token": token, "user": user}})
}

// Me godoc
// @Summary Current user
// @Description Get current authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "unauthenticated"})
		return
	}
	id, ok := uid.(uint)
	if !ok {
		// sometimes numeric types come as float64 from JSON; handle int conversion
		if if64, ok := uid.(int64); ok {
			id = uint(if64)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "invalid user id"})
			return
		}
	}
	u, err := h.svc.Me(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": u})
}
