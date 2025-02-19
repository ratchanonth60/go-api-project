package response

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
type RegisterResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
