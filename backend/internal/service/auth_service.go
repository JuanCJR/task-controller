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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	hasValidCredentials := utils.CheckPasswordHash(password, user.Password)

	if !hasValidCredentials {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, s.cfg.JwtSecret, s.cfg.TokenExpiration)
	if err != nil {
		return "", err
	}

	return token, nil
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
