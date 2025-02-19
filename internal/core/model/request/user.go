package request

import (
	"project-api/internal/core/entity"
)

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

func (r *UserRequest) ToEntity() *entity.User {
	return &entity.User{
		UserName:  r.Username,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  r.Password,
	}
}
