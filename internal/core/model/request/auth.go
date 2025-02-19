package request

type LoginRequest struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	LoginRequest
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
}

func (r *RegisterRequest) ConfirmPassword() bool {
	return r.Password == r.PasswordConfirm
}
