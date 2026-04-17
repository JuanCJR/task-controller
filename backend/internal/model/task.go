package model

type TaskState string

const (
	TaskStatePending  TaskState = "Pendiente"
	TaskStateAssigned TaskState = "Asignado"
	TaskStateStarted  TaskState = "Iniciado"
	TaskStateOnHold   TaskState = "En espera"
	TaskStateSuccess  TaskState = "Finalizada con exito"
	TaskStateError    TaskState = "Finalizada con error"
)

type Task struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	ExpiresAt   string    `db:"expires_at" json:"expires_at"`
	CreatedAt   string    `db:"created_at" json:"created_at"`
	UpdatedAt   string    `db:"updated_at" json:"updated_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	AssignedTo  string    `db:"assigned_to" json:"assigned_to"`
	TaskState   TaskState `db:"task_state" json:"task_state"`
}
