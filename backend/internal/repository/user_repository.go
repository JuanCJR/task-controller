package repository

import (
	"context"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users where id =$1", id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users where email =$1", email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password, must_change_password, first_name, last_name)
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.MustChangePassword, user.FirstName, user.LastName).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) CreateWithRole(ctx context.Context, user *model.User, roleID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (email, password, must_change_password, first_name, last_name)
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, user.Email, user.Password, user.MustChangePassword, user.FirstName, user.LastName).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
		user.ID, roleID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email=$1, password=$2, must_change_password=$3, first_name=$4, last_name=$5, updated_at=NOW() WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.MustChangePassword, user.FirstName, user.LastName, user.ID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id string, password string, mustChange bool) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET password=$1, must_change_password=$2, updated_at=NOW() WHERE id=$3",
		password, mustChange, id,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
