package model

import "time"

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Module string

const (
	ModuleUser Module = "user"
	ModuleTask Module = "task"
	ModuleRole Module = "role"
)

type Permission struct {
	ID          string    `db:"id" json:"id"`
	Action      Action    `db:"action" json:"action"`
	Description string    `db:"description" json:"description"`
	Module      Module    `db:"module" json:"module"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
