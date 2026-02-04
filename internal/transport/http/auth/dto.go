package auth

type RegisterRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	FullName    string `json:"fullname"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
}
