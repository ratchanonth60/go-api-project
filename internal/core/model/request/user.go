package request

import (
	"project-api/internal/core/entity"
)

type UserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	IsActive     bool   `json:"is_active"`
	ConfirmToken string `json:"confirm_token"`
}

func (r *UserRequest) ToEntity() (*entity.User, error) {
	if r.Email == "" {
		return nil, ErrEmailRequired // Use custom error type
	}
	if !isValidEmail(r.Email) {
		return nil, ErrInvalidEmail // Use custom error type
	}
	if r.Username == "" {
		return nil, ErrUsernameRequired
	}
	if r.Password == "" {
		return nil, ErrPasswordRequired
	}

	return &entity.User{
		UserName:     r.Username,
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Password:     r.Password,
		IsActive:     r.IsActive,
		ConfirmToken: r.ConfirmToken,
	}, nil
}
