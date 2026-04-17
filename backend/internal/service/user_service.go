package service

import (
	"context"
	"fmt"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/JuanCJR/task-controller/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) Create(ctx context.Context, email, password, firstName, lastName, roleName string) (*model.User, error) {
	// No se puede crear usuario con rol Admin
	if roleName == "Admin" {
		return nil, fmt.Errorf("cannot create users with Admin role")
	}

	// Verificar que el rol existe
	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("role '%s' not found", roleName)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password")
	}

	user := &model.User{
		Email:              email,
		Password:           hashedPassword,
		MustChangePassword: true,
		FirstName:          firstName,
		LastName:           lastName,
	}

	err = s.userRepo.CreateWithRole(ctx, user, role.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id, email, firstName, lastName string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if email != "" {
		user.Email = email
	}
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
