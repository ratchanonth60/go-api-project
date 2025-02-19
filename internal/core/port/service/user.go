package service

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/model/request"
)

type IUserService interface {
	Create(ctx context.Context, user *request.UserRequest) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUserName(ctx context.Context, username string) (*entity.User, error)
}
