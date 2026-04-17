package service

import (
	"context"
	"fmt"

	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/JuanCJR/task-controller/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      config.AuthConfig
}

func NewAuthService(userRepo *repository.UserRepository, cfg config.AuthConfig) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type LoginResponse struct {
	Token              string `json:"token"`
	MustChangePassword bool   `json:"must_change_password"`
	UserID             string `json:"user_id"`
	Email              string `json:"email"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, s.cfg.JwtSecret, s.cfg.TokenExpiration)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              token,
		MustChangePassword: user.MustChangePassword,
		UserID:             user.ID,
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return fmt.Errorf("invalid credentials")
	}
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password")
	}
	user.Password = newPasswordHash
	user.MustChangePassword = false
	return s.userRepo.UpdatePassword(ctx, user.ID, user.Password, user.MustChangePassword)
}
