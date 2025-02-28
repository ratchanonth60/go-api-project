package repository

import (
	"context"
	"strings"

	"project-api/internal/core/entity"
	"project-api/internal/core/port/repository"
	"project-api/internal/infra/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, entity *entity.User) error {
	return u.db.WithContext(ctx).Create(entity).Error
}

func (u *UserRepository) GetById(ctx context.Context, id uint) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("email = ?", strings.ToLower(email)).Or("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("user_name = ? AND is_active = true", name).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("confirm_token = ?", token).First(&user).Error
	if err != nil {
		logger.Error("Failed to find user by token", zap.String("token", token), zap.Error(err))
		return nil, err
	}
	logger.Info("User found by token",
		zap.String("token", token),
		zap.String("email", user.Email),
		zap.Bool("is_verified", user.IsActive))
	return &user, nil
}

func (u *UserRepository) Update(ctx context.Context, entity *entity.User) error {
	return u.db.WithContext(ctx).Save(entity).Error
}

func (u *UserRepository) FindByResetToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("reset_password_token = ?", token).First(&user).Error
	if err != nil {
		logger.Error("Failed to find user by reset token", zap.String("token", token), zap.Error(err))
		return nil, err
	}
	return &user, nil
}
