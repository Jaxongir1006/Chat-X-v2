package domain

import "time"

const (
	UserRoleAdmin     = "admin"
	UserRoleUser      = "user"
	UserRoleSuperuser = "superuser"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password_hash"`
	Verified  bool      `json:"verified"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID           uint64    `json:"id"`
	FullName     string    `json:"fullname"`
	Address      string    `json:"address"`
	ProfileImage string    `json:"profile_image_link"`
	UserID       uint64    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserSession struct {
	ID              uint64     `json:"id"`
	UserID          uint64     `json:"user_id"`
	AccessToken     string     `json:"access_token"`
	AccessTokenExp  time.Time  `json:"access_token_expires_at"`
	RefreshToken    string     `json:"refresh_token"`
	RefreshTokenExp time.Time  `json:"refresh_token_expires_at"`
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	Device          string     `json:"device"`
	RevokedAt       *time.Time `json:"revoked_at"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	User            *User      `json:"user,omitempty"`
}
