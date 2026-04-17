package database

import (
	"fmt"

	"github.com/JuanCJR/task-controller/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewConnection(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("error connection to database: %w", err)
	}

	return db, nil

}
