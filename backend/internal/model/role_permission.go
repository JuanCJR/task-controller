package model

import "time"

type RolePermission struct {
	RoleID       string    `db:"role_id" json:"role_id"`
	PermissionID string    `db:"permission_id" json:"permission_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
