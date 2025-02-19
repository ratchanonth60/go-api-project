package service

import (
	"context"
	"errors"

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
	userEntity, err := user.ToEntity()
	if err != nil {
		return err
	}
	if userEntity == nil {
		return ErrCreateUser
	}
	if err := u.repo.Create(ctx, userEntity); err != nil {
		return wrapError(ErrCreateUser, err) // Wrap repository errors
	}

	return nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, wrapError(errors.New("failed to get user by email"), err) // Wrap
	}
	return user, nil
}

func (u *UserService) GetUserByUserName(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.repo.GetUserByName(ctx, username)
	if err != nil {
		return nil, wrapError(errors.New("failed to get user by username"), err) // Wrap
	}
	return user, nil
}
