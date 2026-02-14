package authUsecase

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
}

type VerifyUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  int    `json:"code" binding:"required"`
}

type VerifyUserResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessTokenExp  string `json:"access_token_ttl"`
	RefreshTokenExp string `json:"refresh_token_ttl"`
	Device          string `json:"device"`
	UserEmail       string `json:"user_email"`
	IpAddress       string `json:"ip_address"`
}

type SessionMeta struct {
	IP        string
	UserAgent string
	Device    string
}

type LoginRequest struct {
	LoginInput string `json:"login_input" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessTokenExp  string `json:"access_token_ttl"`
	RefreshTokenExp string `json:"refresh_token_ttl"`
	Device          string `json:"device"`
	UserEmail       string `json:"user_email"`
	IpAddress       string `json:"ip_address"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessTokenExp  string `json:"access_token_ttl"`
	RefreshTokenExp string `json:"refresh_token_ttl"`
	Device          string `json:"device"`
	IpAddress       string `json:"ip_address"`
}
