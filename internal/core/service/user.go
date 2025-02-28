package service

import (
	"context"
	"errors"
	"fmt"

	"project-api/internal/core/entity"
	In "project-api/internal/core/port/repository"
	"project-api/internal/infra/logger"
	"project-api/internal/infra/redis"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  In.IUserRepository
	redis *redis.RedisClient
}

func NewUserService(repo In.IUserRepository) *UserService {
	return &UserService{
		repo:  repo,
		redis: nil,
	}
}

func (u *UserService) Create(ctx context.Context, user *entity.User) error {

	if err := u.repo.Create(ctx, user); err != nil {
		return wrapError(ErrCreateUser, err) // Wrap repository errors
	}
	// u.invalidateCache(ctx, user)
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

func (u *UserService) ResendConfirmationEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to find user for resend confirmation", zap.String("email", email), zap.Error(err))
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.IsActive {
		logger.Info("Email already verified, no need to resend", zap.String("email", user.Email))
		return nil, errors.New("email already verified")
	}

	// สร้าง token ใหม่
	newToken := uuid.New().String()
	user.ConfirmToken = newToken
	if err := u.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user with new confirmation token", zap.String("email", email), zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("New confirmation token generated", zap.String("email", user.Email), zap.String("token", newToken))
	return user, nil
}

func (u *UserService) ResetPassword(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to find user for reset password", zap.String("email", email), zap.Error(err))
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// สร้าง reset password token
	resetToken := uuid.New().String()
	user.ResetPasswordToken = resetToken
	if err := u.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user with reset password token", zap.String("email", email), zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("Reset password token generated", zap.String("email", user.Email), zap.String("token", resetToken))
	return user, nil
}

func (u *UserService) ConfirmResetPassword(ctx context.Context, token string, newPassword string) error {
	// ค้นหา user จาก reset token
	user, err := u.repo.FindByResetToken(ctx, token)
	if err != nil {
		logger.Error("Invalid or expired reset token", zap.String("token", token), zap.Error(err))
		return fmt.Errorf("invalid or expired reset token: %w", err)
	}

	// เข้ารหัสรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash new password", zap.Error(err))
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// อัปเดต user: ลบ reset token และตั้งรหัสผ่านใหม่
	user.ResetPasswordToken = ""
	user.Password = string(hashedPassword)
	if err := u.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user with new password", zap.String("email", user.Email), zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("Password reset successfully", zap.String("email", user.Email))
	return nil
}

func (u *UserService) invalidateCache(ctx context.Context, user *entity.User) {
	keys := []string{
		fmt.Sprintf("user:id:%d", user.ID),
		fmt.Sprintf("user:email:%s", user.Email),
		fmt.Sprintf("user:name:%s", user.UserName),
	}
	for _, key := range keys {
		if err := u.redis.DeleteFromCache(ctx, key); err != nil {
			logger.Warn("Failed to delete cache", zap.String("key", key), zap.Error(err))
		}
	}
}
