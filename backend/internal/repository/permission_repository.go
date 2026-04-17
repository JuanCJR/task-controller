package repository

import (
	"context"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/jmoiron/sqlx"
)

type PermissionRepository struct {
	db *sqlx.DB
}

func NewPermissionRepository(db *sqlx.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) UserHasPermission(ctx context.Context, userID string, action model.Action, module model.Module) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) FROM permissions p
        INNER JOIN role_permissions rp ON rp.permission_id = p.id
        INNER JOIN user_roles ur ON ur.role_id = rp.role_id
        WHERE ur.user_id = $1 AND p.action = $2 AND p.module = $3
    `
	err := r.db.GetContext(ctx, &count, query, userID, action, module)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PermissionRepository) NeedToChangePassword(ctx context.Context, userID string) (bool, error) {
	var mustChange bool
	query := `SELECT must_change_password FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &mustChange, query, userID)
	if err != nil {
		return false, err
	}
	return mustChange, nil
}
