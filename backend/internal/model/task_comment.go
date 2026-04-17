package model

type TaskComment struct {
	ID        string `db:"id" json:"id"`
	TaskID    string `db:"task_id" json:"task_id"`
	UserID    string `db:"user_id" json:"user_id"`
	Comment   string `db:"comment" json:"comment"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
