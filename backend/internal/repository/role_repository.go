package repository

import (
	"context"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := r.db.GetContext(ctx, &role, "SELECT * FROM roles WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetUserRole(ctx context.Context, userID string) (*model.Role, error) {
	var role model.Role
	query := `
		SELECT r.* FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`
	err := r.db.GetContext(ctx, &role, query, userID)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
		userID, roleID,
	)
	return err
}
