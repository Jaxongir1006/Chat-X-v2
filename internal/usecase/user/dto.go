package userUsecase

import "time"


type UserResponse struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Username string `json:"username"`
	Verified bool `json:"verified"`
	Phone string `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Profile   UserProfileResponse `json:"profile"`
}

type UserProfileResponse struct {
	FullName     string    `json:"fullname"`
	Address      string    `json:"address"`
	Bio          string    `json:"bio"`
	ProfileImage string    `json:"profile_image_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FullName 	*string `json:"fullname"`
	Address  	*string `json:"address"`
	ProfileImage *string `json:"profile_image_key"`
	Bio			*string `json:"bio"`
}