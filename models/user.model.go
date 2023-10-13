package models

import (
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uuid.UUID      `json:"id"                             gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email             string         `json:"email"                          gorm:"uniqueIndex;not null"`
	Password          string         `json:"-"                              gorm:"not null"`
	Role              Role           `json:"role"                           gorm:"type:varchar(255);not null"`
	Verified          bool           `json:"-"                              gorm:"not null"`
	VerificationToken string         `json:"verification_token,omitempty"`
	ResetToken        string         `json:"reset_token,omitempty"`
	CreatedAt         time.Time      `json:"-"                              gorm:"not null"`
	UpdatedAt         time.Time      `json:"-"                              gorm:"not null"`
	ResetedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"                              gorm:"index"`
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

type UserPasswordConfirmation struct {
	Password        string `json:"password"         binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

func (u *User) Response() User {

	// User response for the client
	user := User{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}

	// Add the verification code if the app is in debug mode
	if config.Config.App.Debug {
		if len(u.VerificationToken) > 0 {
			user.VerificationToken = u.VerificationToken
		}
		if len(u.ResetToken) > 0 {
			user.ResetToken = u.ResetToken
		}
	}

	return user
}
