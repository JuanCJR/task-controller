package model

import "time"

type User struct {
	ID                 string    `db:"id" json:"id"`
	Email              string    `db:"email" json:"email"`
	Password           string    `db:"password"`
	MustChangePassword bool      `db:"must_change_password" json:"must_change_password"`
	FirstName          string    `db:"first_name" json:"first_name"`
	LastName           string    `db:"last_name" json:"last_name"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}
