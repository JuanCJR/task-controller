package handler

import (
	"net/http"

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
