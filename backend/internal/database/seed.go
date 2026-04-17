package database

import (
	"context"
	"log"

	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/pkg/utils"
	"github.com/jmoiron/sqlx"
)

func Seed(db *sqlx.DB, cfg config.AppConfig) {
	ctx := context.Background()

	// Verificar si ya existe data
	var userCount int
	err := db.GetContext(ctx, &userCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		log.Printf("Seed: error checking users table: %v", err)
		return
	}
	if userCount > 0 {
		log.Println("Seed: database already has data, skipping...")
		return
	}

	log.Println("Seed: seeding database...")

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Seed: error starting transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// ==================== ROLES ====================
	var adminRoleID, ejecutorRoleID, auditorRoleID string

	err = tx.QueryRowContext(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id", "Admin",
	).Scan(&adminRoleID)
	if err != nil {
		log.Printf("Seed: error creating Admin role: %v", err)
		return
	}

	err = tx.QueryRowContext(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id", "Ejecutor",
	).Scan(&ejecutorRoleID)
	if err != nil {
		log.Printf("Seed: error creating Ejecutor role: %v", err)
		return
	}

	err = tx.QueryRowContext(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id", "Auditor",
	).Scan(&auditorRoleID)
	if err != nil {
		log.Printf("Seed: error creating Auditor role: %v", err)
		return
	}

	log.Println("Seed: roles created")

	// ==================== PERMISSIONS ====================
	type permDef struct {
		Action      string
		Module      string
		Description string
	}

	permissions := []permDef{
		// User module
		{"create", "user", "Create users"},
		{"read", "user", "Read users"},
		{"update", "user", "Update users"},
		{"delete", "user", "Delete users"},
		// Task module
		{"create", "task", "Create tasks"},
		{"read", "task", "Read tasks"},
		{"update", "task", "Update tasks"},
		{"delete", "task", "Delete tasks"},
		// Role module
		{"create", "role", "Create roles"},
		{"read", "role", "Read roles"},
		{"update", "role", "Update roles"},
		{"delete", "role", "Delete roles"},
	}

	permIDs := make(map[string]string) // key: "action:module" → value: permission ID

	for _, p := range permissions {
		var id string
		err = tx.QueryRowContext(ctx,
			"INSERT INTO permissions (action, module, description) VALUES ($1, $2, $3) RETURNING id",
			p.Action, p.Module, p.Description,
		).Scan(&id)
		if err != nil {
			log.Printf("Seed: error creating permission %s:%s: %v", p.Action, p.Module, err)
			return
		}
		permIDs[p.Action+":"+p.Module] = id
	}

	log.Println("Seed: permissions created")

	// ==================== ROLE PERMISSIONS ====================

	// Admin: todos los permisos
	for _, id := range permIDs {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)",
			adminRoleID, id,
		)
		if err != nil {
			log.Printf("Seed: error assigning permission to Admin: %v", err)
			return
		}
	}

	// Ejecutor: read y update de tasks
	ejecutorPerms := []string{"read:task", "update:task"}
	for _, key := range ejecutorPerms {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)",
			ejecutorRoleID, permIDs[key],
		)
		if err != nil {
			log.Printf("Seed: error assigning permission to Ejecutor: %v", err)
			return
		}
	}

	// Auditor: read de tasks
	auditorPerms := []string{"read:task"}
	for _, key := range auditorPerms {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)",
			auditorRoleID, permIDs[key],
		)
		if err != nil {
			log.Printf("Seed: error assigning permission to Auditor: %v", err)
			return
		}
	}

	log.Println("Seed: role permissions assigned")

	// ==================== ADMIN USER ====================
	hashedPassword, err := utils.HashPassword(cfg.DefaultAdminPassword)
	if err != nil {
		log.Printf("Seed: error hashing password: %v", err)
		return
	}

	var adminUserID string
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (email, password, must_change_password, first_name, last_name)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		cfg.DefaultAdminEmail, hashedPassword, true, "Admin", "System",
	).Scan(&adminUserID)
	if err != nil {
		log.Printf("Seed: error creating admin user: %v", err)
		return
	}

	// Assign Admin role to admin user
	_, err = tx.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
		adminUserID, adminRoleID,
	)
	if err != nil {
		log.Printf("Seed: error assigning Admin role to user: %v", err)
		return
	}

	// ==================== COMMIT ====================
	err = tx.Commit()
	if err != nil {
		log.Printf("Seed: error committing transaction: %v", err)
		return
	}

	log.Println("Seed: database seeded successfully")
	log.Println("Seed: Admin user created:")
	log.Printf("  Email: %s", cfg.DefaultAdminEmail)
	log.Printf("  Password: %s", cfg.DefaultAdminPassword)
	log.Println("  (must change password on first login)")
}
