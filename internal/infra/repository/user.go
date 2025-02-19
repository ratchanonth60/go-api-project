package repository

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/port/repository"

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

func (u *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return u.db.WithContext(ctx).Create(&user).Error
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("user_name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
