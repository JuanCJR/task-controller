package handler

import (
	"net/http"

	"github.com/JuanCJR/task-controller/internal/dto"
	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService    *service.UserService
	authMiddleware gin.HandlerFunc
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc
}

func NewUserHandler(
	userService *service.UserService,
	authMiddleware gin.HandlerFunc,
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

func (h *UserHandler) RegisterRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.Use(h.authMiddleware)
	{
		users.GET("", h.rbacMiddleware(model.ActionRead, model.ModuleUser), h.GetAll)
		users.POST("", h.rbacMiddleware(model.ActionCreate, model.ModuleUser), h.Create)
		users.PUT("/:id", h.rbacMiddleware(model.ActionUpdate, model.ModuleUser), h.Update)
		users.DELETE("/:id", h.rbacMiddleware(model.ActionDelete, model.ModuleUser), h.Delete)
	}
}

// GetAll godoc
// @Summary      Get all users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  model.User
// @Failure      500  {object} map[string]string
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Param        id path string true "User ID"
// @Security     BearerAuth
// @Success      200  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user with a role (Ejecutor or Auditor only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateUserRequest true "User data"
// @Security     BearerAuth
// @Success      201  {object} model.User
// @Failure      400  {object} map[string]string
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.RoleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Update godoc
// @Summary      Update user
// @Description  Update user information by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        body body dto.UpdateUserRequest true "User data"
// @Security     BearerAuth
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.Update(c.Request.Context(), id, req.Email, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
