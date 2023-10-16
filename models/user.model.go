package models

import (
	"database/sql"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uuid.UUID      `json:"id"                             gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email             string         `json:"email"                          gorm:"uniqueIndex;not null"`
	Name              string         `json:"name"                           gorm:"not null"`
	FirstName         string         `json:"first_name"                     gorm:"not null"`
	LastName          string         `json:"last_name"                      gorm:"not null"`
	Password          string         `json:"-"                              gorm:"not null"`
	Role              Role           `json:"role"                           gorm:"type:varchar(255);not null"`
	Verified          bool           `json:"-"                              gorm:"not null"`
	VerificationToken string         `json:"verification_token,omitempty"`
	ResetToken        string         `json:"reset_token,omitempty"`
	CreatedAt         time.Time      `json:"-"                              gorm:"not null"`
	UpdatedAt         time.Time      `json:"-"                              gorm:"not null"`
	ResetedAt         sql.NullTime   `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"                              gorm:"index"`
}

type UserSignUp struct {
	Email           string `json:"email"            binding:"required"`
	FirstName       string `json:"first_name"       binding:"required"`
	LastName        string `json:"last_name"        binding:"required"`
	Password        string `json:"password"         binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type UserSignIn struct {
	Email    string `json:"email"     binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserPasswordConfirmation struct {
	Password        string `json:"password"         binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

// SetResetToken sets the reset token
func (u *User) SetResetToken(token string) {
	u.ResetToken = token
	u.ResetedAt = sql.NullTime{Time: time.Now(), Valid: true}
}

// Response returns the user without the password
func (u *User) Response() User {

	// User response for the client
	user := User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}

	// Add the verification code if the app is in debug mode
	if config.Config.App.Debug {
		if len(u.VerificationToken) > 0 {
			user.VerificationToken = u.VerificationToken
		}
		if len(u.ResetToken) > 0 {
			user.SetResetToken(u.ResetToken)
		}
	}

	return user
}
