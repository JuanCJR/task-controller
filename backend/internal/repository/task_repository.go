package repository

import (
	"context"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `INSERT INTO tasks (title, description, expires_at, created_by, assigned_to, task_state)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		task.Title, task.Description, task.ExpiresAt,
		task.CreatedBy, task.AssignedTo, task.TaskState,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	err := r.db.GetContext(ctx, &task, "SELECT * FROM tasks WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.SelectContext(ctx, &tasks, "SELECT * FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) GetByAssignedTo(ctx context.Context, userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.SelectContext(ctx, &tasks,
		"SELECT * FROM tasks WHERE assigned_to = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, expires_at=$3, assigned_to=$4, updated_at=NOW()
	          WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.ExpiresAt, task.AssignedTo, task.ID)
	return err
}

func (r *TaskRepository) UpdateState(ctx context.Context, id string, state model.TaskState) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET task_state=$1, updated_at=NOW() WHERE id=$2",
		state, id)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1", id)
	return err
}
