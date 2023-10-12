package models

import (
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID uuid.UUID `json:"id"                            gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	//Name             string    `json:"name"                          gorm:"type:varchar(255);not null"`
	Email            string    `json:"email"                         gorm:"uniqueIndex;not null"`
	Password         string    `json:"-"                             gorm:"not null"`
	Role             Role      `json:"role"                          gorm:"type:varchar(255);not null"`
	VerificationCode string    `json:"verification_code,omitempty"`
	Verified         bool      `json:"-"                             gorm:"not null"`
	CreatedAt        time.Time `json:"-"                             gorm:"not null"`
	UpdatedAt        time.Time `json:"-"                             gorm:"not null"`
}

type UserSignUp struct {
	Name            string `json:"name"             binding:"required"`
	Email           string `json:"email"            binding:"required"`
	Password        string `json:"password"         binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type UserSignIn struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

func (u *User) Response() User {

	// User response for the client
	user := User{
		ID: u.ID,
		//Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}

	// Add the verification code if the app is in debug mode
	if config.Config.App.Debug {
		if len(u.VerificationCode) > 0 {
			user.VerificationCode = u.VerificationCode
		}
	}

	return user
}
