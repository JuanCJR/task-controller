package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
)

type TaskCommentService struct {
	commentRepo *repository.TaskCommentRepository
	taskRepo    *repository.TaskRepository
	roleRepo    *repository.RoleRepository
}

func NewTaskCommentService(
	commentRepo *repository.TaskCommentRepository,
	taskRepo *repository.TaskRepository,
	roleRepo *repository.RoleRepository,
) *TaskCommentService {
	return &TaskCommentService{commentRepo: commentRepo, taskRepo: taskRepo, roleRepo: roleRepo}
}

// Create - Ejecutor puede agregar comentarios sobre tareas vencidas.
func (s *TaskCommentService) Create(ctx context.Context, taskID, userID, comment string) (*model.TaskComment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	// Verificar rol del usuario
	role, err := s.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user role not found")
	}

	// Ejecutor solo puede comentar tareas asignadas a él y que estén vencidas
	if role.Name == "Ejecutor" {
		if task.AssignedTo != userID {
			return nil, fmt.Errorf("you can only comment on tasks assigned to you")
		}
		if task.ExpiresAt.IsZero() {
			return nil, fmt.Errorf("task has no expiration date")
		}
		if !time.Now().After(task.ExpiresAt) {
			return nil, fmt.Errorf("ejecutor can only comment on expired tasks")
		}
	}

	tc := &model.TaskComment{
		TaskID:  taskID,
		UserID:  userID,
		Comment: comment,
	}
	err = s.commentRepo.Create(ctx, tc)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

// GetByTaskID - Obtener comentarios de una tarea.
func (s *TaskCommentService) GetByTaskID(ctx context.Context, taskID string) ([]model.TaskComment, error) {
	return s.commentRepo.GetByTaskID(ctx, taskID)
}
