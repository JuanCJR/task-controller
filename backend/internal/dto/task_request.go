package dto

import "time"

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	AssignedTo  string    `json:"assigned_to" binding:"required,uuid"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title" binding:"omitempty"`
	Description string     `json:"description" binding:"omitempty"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"omitempty"`
	AssignedTo  string     `json:"assigned_to" binding:"omitempty,uuid"`
}

type UpdateTaskStateRequest struct {
	TaskState string `json:"task_state" binding:"required"`
}
