package user

import "time"

type UserResponse struct {
	ID        uint64              `json:"id"`
	Email     string              `json:"email"`
	Role      string              `json:"role"`
	Username  string              `json:"username"`
	Verified  bool                `json:"verified"`
	Phone     string              `json:"phone"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Profile   UserProfileResponse `json:"profile"`
}

type UserProfileResponse struct {
	FullName     string                `json:"fullname"`
	Address      string                `json:"address"`
	Bio          string                `json:"bio"`
	Media        []UserProfileMediaDTO `json:"media,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

type UserProfileMediaDTO struct {
	ID           uint64 `json:"id"`
	MediaKey     string `json:"media_key"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int    `json:"display_order"`
}

type UpdateProfileRequest struct {
	FullName     *string `json:"fullname"`
	Address      *string `json:"address"`
	Bio          *string `json:"bio"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type AddProfileMediaRequest struct {
	MediaKey  string `json:"media_key" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}
