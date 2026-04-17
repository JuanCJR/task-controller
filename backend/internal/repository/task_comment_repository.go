package repository

import (
	"context"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/jmoiron/sqlx"
)

type TaskCommentRepository struct {
	db *sqlx.DB
}

func NewTaskCommentRepository(db *sqlx.DB) *TaskCommentRepository {
	return &TaskCommentRepository{db: db}
}

func (r *TaskCommentRepository) Create(ctx context.Context, comment *model.TaskComment) error {
	query := `INSERT INTO task_comments (task_id, user_id, comment)
	          VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		comment.TaskID, comment.UserID, comment.Comment,
	).Scan(&comment.ID, &comment.CreatedAt)
}

func (r *TaskCommentRepository) GetByTaskID(ctx context.Context, taskID string) ([]model.TaskComment, error) {
	var comments []model.TaskComment
	err := r.db.SelectContext(ctx, &comments,
		"SELECT * FROM task_comments WHERE task_id = $1 ORDER BY created_at DESC", taskID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
