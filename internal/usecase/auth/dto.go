package authUsecase

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
}
