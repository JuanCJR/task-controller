package model

import "time"

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
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	AssignedTo  string    `db:"assigned_to" json:"assigned_to"`
	TaskState   TaskState `db:"task_state" json:"task_state"`
}
