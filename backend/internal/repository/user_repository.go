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
