package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// LoginRequest represents the data structure for login requests
type LoginRequest struct {
	UserName string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=100"`
}

// RegisterRequest represents the data structure for registration requests
type RegisterRequest struct {
	LoginRequest
	FirstName       string `json:"first_name" validate:"required,min=2,max=50"`
	LastName        string `json:"last_name" validate:"required,min=2,max=50"`
	Email           string `json:"email" validate:"required,email"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm" validate:"required"`
}

// ConfirmPassword checks if the password and confirmation match
func (r *RegisterRequest) ConfirmPassword() bool {
	return r.Password == r.PasswordConfirm
}

// Validate validates the LoginRequest struct
func (r *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate validates the RegisterRequest struct, including password confirmation
func (r *RegisterRequest) Validate() error {
	validate := validator.New()

	// Validate the struct fields
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Custom validation for password confirmation
	if !r.ConfirmPassword() {
		return fmt.Errorf("password and password confirmation do not match")
	}

	return nil
}
