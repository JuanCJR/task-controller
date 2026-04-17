package handler

import (
	"net/http"

	"github.com/JuanCJR/task-controller/internal/dto"
	"github.com/JuanCJR/task-controller/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService    *service.AuthService
	authMiddleware gin.HandlerFunc
}

func NewAuthHandler(authService *service.AuthService, authMiddleware gin.HandlerFunc) *AuthHandler {
	return &AuthHandler{authService: authService, authMiddleware: authMiddleware}
}

func (h *AuthHandler) RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)

	}
	auth.Use(h.authMiddleware)
	{
		auth.PUT("/change-password", h.ChangePassword)
	}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login credentials"
// @Success      200  {object} service.LoginResponse
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)

}

// Logout godoc
// @Summary      Logout
// @Description  Logout user
// @Tags         auth
// @Produce      json
// @Success      200  {object} map[string]string
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change user password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.ChangePasswordRequest true "Password data"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Security     BearerAuth
// @Router       /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.authService.ChangePassword(c.Request.Context(), userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
