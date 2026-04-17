package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
)

// Transiciones de estado permitidas para Ejecutor
var allowedTransitions = map[model.TaskState][]model.TaskState{
	model.TaskStateAssigned: {model.TaskStateStarted},
	model.TaskStateStarted:  {model.TaskStateOnHold, model.TaskStateSuccess, model.TaskStateError},
	model.TaskStateOnHold:   {model.TaskStateStarted, model.TaskStateSuccess, model.TaskStateError},
}

// ValidTaskStates contiene todos los estados válidos
var validTaskStates = map[model.TaskState]bool{
	model.TaskStatePending:  true,
	model.TaskStateAssigned: true,
	model.TaskStateStarted:  true,
	model.TaskStateOnHold:   true,
	model.TaskStateSuccess:  true,
	model.TaskStateError:    true,
}

type TaskService struct {
	taskRepo *repository.TaskRepository
	roleRepo *repository.RoleRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, roleRepo *repository.RoleRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, roleRepo: roleRepo}
}

// Create - Solo Admin puede crear tareas. assigned_to debe ser Ejecutor.
func (s *TaskService) Create(ctx context.Context, task *model.Task) error {
	assignedRole, err := s.roleRepo.GetUserRole(ctx, task.AssignedTo)
	if err != nil {
		return fmt.Errorf("assigned user not found")
	}
	if assignedRole.Name != "Ejecutor" {
		return fmt.Errorf("tasks can only be assigned to users with Ejecutor role")
	}

	task.TaskState = model.TaskStateAssigned
	return s.taskRepo.Create(ctx, task)
}

// GetAll - Admin y Auditor ven todas las tareas.
func (s *TaskService) GetAll(ctx context.Context) ([]model.Task, error) {
	return s.taskRepo.GetAll(ctx)
}

// GetByID - Obtener tarea por ID.
func (s *TaskService) GetByID(ctx context.Context, id string) (*model.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// GetByAssignedTo - Ejecutor solo ve sus tareas.
func (s *TaskService) GetByAssignedTo(ctx context.Context, userID string) ([]model.Task, error) {
	return s.taskRepo.GetByAssignedTo(ctx, userID)
}

// Update - Admin actualiza tareas. Solo si estado es "Asignado".
func (s *TaskService) Update(ctx context.Context, taskID string, title, description string, expiresAt *time.Time, assignedTo string) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	if task.TaskState != model.TaskStateAssigned {
		return fmt.Errorf("task can only be updated when in 'Asignado' state")
	}

	if assignedTo != "" && assignedTo != task.AssignedTo {
		assignedRole, err := s.roleRepo.GetUserRole(ctx, assignedTo)
		if err != nil {
			return fmt.Errorf("assigned user not found")
		}
		if assignedRole.Name != "Ejecutor" {
			return fmt.Errorf("tasks can only be assigned to users with Ejecutor role")
		}
		task.AssignedTo = assignedTo
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if expiresAt != nil {
		task.ExpiresAt = *expiresAt
	}

	return s.taskRepo.Update(ctx, task)
}

// UpdateState - Ejecutor actualiza estado de sus tareas. No si está vencida. Valida transiciones.
func (s *TaskService) UpdateState(ctx context.Context, taskID, userID string, newState model.TaskState) error {
	if !validTaskStates[newState] {
		return fmt.Errorf("invalid task state: %s", newState)
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	if task.AssignedTo != userID {
		return fmt.Errorf("you can only update tasks assigned to you")
	}

	// Verificar si la tarea está vencida
	if !task.ExpiresAt.IsZero() && time.Now().After(task.ExpiresAt) {
		return fmt.Errorf("cannot update state of an expired task")
	}

	// Validar transición de estado
	allowed, exists := allowedTransitions[task.TaskState]
	if !exists {
		return fmt.Errorf("task in state '%s' cannot be updated", task.TaskState)
	}
	isAllowed := false
	for _, s := range allowed {
		if s == newState {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("cannot transition from '%s' to '%s'", task.TaskState, newState)
	}

	return s.taskRepo.UpdateState(ctx, taskID, newState)
}

// Delete - Admin elimina tareas. Solo si estado es "Asignado".
func (s *TaskService) Delete(ctx context.Context, taskID string) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	if task.TaskState != model.TaskStateAssigned {
		return fmt.Errorf("task can only be deleted when in 'Asignado' state")
	}

	return s.taskRepo.Delete(ctx, taskID)
}
