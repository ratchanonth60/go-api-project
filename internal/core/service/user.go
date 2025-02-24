package service

import (
	"context"
	"errors"

	"project-api/internal/core/entity"
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
