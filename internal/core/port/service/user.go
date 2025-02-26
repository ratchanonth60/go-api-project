package service

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/port/utils"
)

type IUserService interface {
	utils.BaseInterface[entity.User]
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByName(ctx context.Context, username string) (*entity.User, error)
	ConfirmEmail(ctx context.Context, token string) error
}
