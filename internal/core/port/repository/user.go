package repository

import (
	"context"

	"project-api/internal/core/entity"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
}
