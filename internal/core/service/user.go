package service

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/model/request"
	In "project-api/internal/core/port/repository"
)

type UserService struct {
	repo In.IUserRepository
}

func NewUserService(repo In.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Create(ctx context.Context, user *request.UserRequest) error {
	return u.repo.Create(ctx, user.ToEntity())
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *UserService) GetUserByUserName(ctx context.Context, username string) (*entity.User, error) {
	return u.repo.GetUserByName(ctx, username)
}
