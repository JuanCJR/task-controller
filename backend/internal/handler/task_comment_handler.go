package handler

import (
	"net/http"

	"github.com/JuanCJR/task-controller/internal/dto"
	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskCommentHandler struct {
	commentService *service.TaskCommentService
	authMiddleware gin.HandlerFunc
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc
}

func NewTaskCommentHandler(
	commentService *service.TaskCommentService,
	authMiddleware gin.HandlerFunc,
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc,
) *TaskCommentHandler {
	return &TaskCommentHandler{
		commentService: commentService,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

func (h *TaskCommentHandler) RegisterRoutes(api *gin.RouterGroup) {
	comments := api.Group("/tasks/:id/comments")
	comments.Use(h.authMiddleware)
	{
		comments.GET("", h.rbacMiddleware(model.ActionRead, model.ModuleTask), h.GetByTaskID)
		comments.POST("", h.rbacMiddleware(model.ActionUpdate, model.ModuleTask), h.Create)
	}
}

// Create godoc
// @Summary      Create task comment
// @Description  Add a comment to a task. Ejecutor can only comment on expired tasks assigned to them
// @Tags         task-comments
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Param        body body dto.CreateTaskCommentRequest true "Comment data"
// @Security     BearerAuth
// @Success      201  {object} model.TaskComment
// @Failure      400  {object} map[string]string
// @Router       /tasks/{id}/comments [post]
func (h *TaskCommentHandler) Create(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("userID")

	var req dto.CreateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentService.Create(c.Request.Context(), taskID, userID.(string), req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetByTaskID godoc
// @Summary      Get task comments
// @Description  Get all comments for a specific task
// @Tags         task-comments
// @Produce      json
// @Param        id path string true "Task ID"
// @Security     BearerAuth
// @Success      200  {array}  model.TaskComment
// @Failure      500  {object} map[string]string
// @Router       /tasks/{id}/comments [get]
func (h *TaskCommentHandler) GetByTaskID(c *gin.Context) {
	taskID := c.Param("id")

	comments, err := h.commentService.GetByTaskID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
