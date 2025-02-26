package service

import (
	"context"
	"errors"
	"fmt"

	"project-api/internal/core/entity"
	In "project-api/internal/core/port/repository"
	"project-api/internal/infra/logger"

	"go.uber.org/zap"
)

type UserService struct {
	repo In.IUserRepository
}

func NewUserService(repo In.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Create(ctx context.Context, user *entity.User) error {

	if err := u.repo.Create(ctx, user); err != nil {
		return wrapError(ErrCreateUser, err) // Wrap repository errors
	}

	return nil
}

func (u *UserService) GetById(ctx context.Context, id uint) (*entity.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, wrapError(errors.New("failed to get user by email"), err) // Wrap
	}
	return user, nil
}

func (u *UserService) GetUserByName(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.repo.GetUserByName(ctx, username)
	if err != nil {
		return nil, wrapError(errors.New("failed to get user by username"), err) // Wrap
	}
	return user, nil
}

func (s *UserService) ConfirmEmail(ctx context.Context, token string) error {
	// ค้นหา user จาก confirm_token
	user, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		logger.Error("Invalid or expired token", zap.String("token", token), zap.Error(err))
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	// ตรวจสอบว่า email ถูกยืนยันหรือยัง
	if user.IsActive {
		logger.Info("Email already verified", zap.String("email", user.Email))
		return fmt.Errorf("email already verified")
	}

	// อัปเดตสถานะเป็น verified และลบ token
	user.IsActive = true
	user.ConfirmToken = "" // ลบ token หลังยืนยัน
	if err := s.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user verification status", zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("Email confirmed successfully", zap.String("email", user.Email))
	return nil
}

func (u *UserService) Update(ctx context.Context, entity *entity.User) error {
	return u.repo.Update(ctx, entity)
}
