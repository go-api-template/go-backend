package models

//
import (
	"database/sql"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User model
//
//	@description	User model
type User struct {
	// User information
	ID        uuid.UUID `json:"id"         gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string    `json:"email"      gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"       gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name"  gorm:"not null"`
	Password  string    `json:"-"          gorm:"not null"`
	Role      Role      `json:"role"       gorm:"type:varchar(255);not null"`

	// User status
	Verified          bool   `json:"-"     gorm:"not null"`
	VerificationToken string `json:"verification_token,omitempty"`

	// Password reset
	ResetToken  string       `json:"reset_token,omitempty"`
	LastResetAt sql.NullTime `json:"-"`

	// Timestamps
	CreatedAt time.Time      `json:"-"      gorm:"not null"`
	UpdatedAt time.Time      `json:"-"      gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"-"      gorm:"index"`
}

// UserSignUp model
//
//	@description	User sign up model used for registration
type UserSignUp struct {
	Email                string `json:"email"                 binding:"required"`
	Password             string `json:"password"              binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

// UserSignIn model
//
//	@description	User sign in model used for authentication
type UserSignIn struct {
	Email    string `json:"email"     binding:"required" example:"my-email@gmail.com"`
	Password string `json:"password"  binding:"required" example:"strong-password"`
}

// UserEmail model
//
//	@description	User email model used for password reset
type UserEmail struct {
	Email string `json:"email" binding:"required" example:"my-email@gmail.com"`
}

// USerToken model
//
//	@description	Token used for refresh
type UserToken struct {
	Token string `json:"token" binding:"required"`
}

// UserPasswordConfirmation model
//
//	@description	User password confirmation model used for password reset
type UserPasswordConfirmation struct {
	Password             string `json:"password" binding:"required,min=8" example:"strong-password"`
	PasswordConfirmation string `json:"password_confirmation"  binding:"required" example:"strong-password"`
}

// SetResetToken sets the reset token
func (u *User) SetResetToken(token string) {
	u.ResetToken = token
	u.LastResetAt = sql.NullTime{Time: time.Now(), Valid: true}
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
