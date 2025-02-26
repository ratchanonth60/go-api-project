package repository

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/port/utils"
)

type IUserRepository interface {
	utils.BaseInterface[entity.User]
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
	FindByToken(ctx context.Context, token string) (*entity.User, error)
}
