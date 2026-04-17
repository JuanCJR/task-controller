package dto

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleName  string `json:"role_name" binding:"required,oneof=Ejecutor Auditor"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name" binding:"omitempty"`
	LastName  string `json:"last_name" binding:"omitempty"`
}
