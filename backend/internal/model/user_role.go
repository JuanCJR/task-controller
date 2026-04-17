package model

type UserRole struct {
	UserID string `db:"user_id" json:"user_id"`
	RoleID string `db:"role_id" json:"role_id"`
}
