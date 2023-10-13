package models

import (
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID uuid.UUID `json:"id"                                                  gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	//Name             string    `json:"name"                          gorm:"type:varchar(255);not null"`
	Email              string         `json:"email"                          gorm:"uniqueIndex;not null"`
	Password           string         `json:"-"                              gorm:"not null"`
	Role               Role           `json:"role"                           gorm:"type:varchar(255);not null"`
	Verified           bool           `json:"-"                              gorm:"not null"`
	VerificationCode   string         `json:"verification_code,omitempty"`
	ResetPasswordToken string         `json:"reset_password_token,omitempty"`
	ResetPasswordAt    time.Time      `json:"-"`
	CreatedAt          time.Time      `json:"-"                              gorm:"not null"`
	UpdatedAt          time.Time      `json:"-"                              gorm:"not null"`
	DeletedAt          gorm.DeletedAt `json:"-"                              gorm:"index"`
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
		if len(u.ResetPasswordToken) > 0 {
			user.ResetPasswordToken = u.ResetPasswordToken
		}
	}

	return user
}
