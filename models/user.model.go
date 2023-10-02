package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"type:varchar(255);not null"`
	Verified  bool      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserSignUp struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type UserSignIn struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
