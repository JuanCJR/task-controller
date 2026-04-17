package handler

import (
	"net/http"

	"github.com/JuanCJR/task-controller/internal/dto"
	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/JuanCJR/task-controller/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService    *service.TaskService
	roleRepo       *repository.RoleRepository
	authMiddleware gin.HandlerFunc
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc
}

func NewTaskHandler(
	taskService *service.TaskService,
	roleRepo *repository.RoleRepository,
	authMiddleware gin.HandlerFunc,
	rbacMiddleware func(action model.Action, module model.Module) gin.HandlerFunc,
) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		roleRepo:       roleRepo,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

func (h *TaskHandler) RegisterRoutes(api *gin.RouterGroup) {
	tasks := api.Group("/tasks")
	tasks.Use(h.authMiddleware)
	{
		tasks.GET("", h.rbacMiddleware(model.ActionRead, model.ModuleTask), h.GetTasks)
		tasks.GET("/:id", h.rbacMiddleware(model.ActionRead, model.ModuleTask), h.GetByID)
		tasks.POST("", h.rbacMiddleware(model.ActionCreate, model.ModuleTask), h.Create)
		tasks.PUT("/:id", h.rbacMiddleware(model.ActionUpdate, model.ModuleTask), h.Update)
		tasks.PATCH("/:id/state", h.rbacMiddleware(model.ActionUpdate, model.ModuleTask), h.UpdateState)
		tasks.DELETE("/:id", h.rbacMiddleware(model.ActionDelete, model.ModuleTask), h.Delete)
	}
}

// GetTasks godoc
// @Summary      Get tasks
// @Description  Get tasks. Admin/Auditor see all tasks, Ejecutor only sees assigned tasks
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  model.Task
// @Failure      500  {object} map[string]string
// @Router       /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	role, err := h.roleRepo.GetUserRole(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user role"})
		return
	}

	if role.Name == "Ejecutor" {
		tasks, err := h.taskService.GetByAssignedTo(c.Request.Context(), userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tasks)
		return
	}

	// Admin y Auditor ven todas
	tasks, err := h.taskService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetByID godoc
// @Summary      Get task by ID
// @Description  Get a single task by its ID. Ejecutor can only view tasks assigned to them
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Security     BearerAuth
// @Success      200  {object} model.Task
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Ejecutor solo puede ver tareas asignadas a él
	role, err := h.roleRepo.GetUserRole(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user role"})
		return
	}
	if role.Name == "Ejecutor" && task.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only view tasks assigned to you"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Create godoc
// @Summary      Create task
// @Description  Create a new task. Only Admin can create tasks. Assigned user must have Ejecutor role
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateTaskRequest true "Task data"
// @Security     BearerAuth
// @Success      201  {object} model.Task
// @Failure      400  {object} map[string]string
// @Router       /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		ExpiresAt:   req.ExpiresAt,
		CreatedBy:   userID.(string),
		AssignedTo:  req.AssignedTo,
	}

	err := h.taskService.Create(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Update godoc
// @Summary      Update task
// @Description  Update a task. Only Admin can update. Task must be in 'Asignado' state
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Param        body body dto.UpdateTaskRequest true "Task data"
// @Security     BearerAuth
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.taskService.Update(c.Request.Context(), id, req.Title, req.Description, req.ExpiresAt, req.AssignedTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// UpdateState godoc
// @Summary      Update task state
// @Description  Update the state of a task. Ejecutor can only update assigned tasks that are not expired
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Param        body body dto.UpdateTaskStateRequest true "New task state"
// @Security     BearerAuth
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /tasks/{id}/state [patch]
func (h *TaskHandler) UpdateState(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var req dto.UpdateTaskStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.taskService.UpdateState(c.Request.Context(), id, userID.(string), model.TaskState(req.TaskState))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task state updated successfully"})
}

// Delete godoc
// @Summary      Delete task
// @Description  Delete a task. Only Admin can delete. Task must be in 'Asignado' state
// @Tags         tasks
// @Param        id path string true "Task ID"
// @Security     BearerAuth
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.taskService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
