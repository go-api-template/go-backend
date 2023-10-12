package services

import (
	"context"
	"errors"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"strings"
)

// AuthService is the interface that the service must implement
// It declares the methods that the service must implement
type AuthService interface {
	UserSignUp(user *models.UserSignUp) (*models.User, error)
}

// AuthServiceImpl is the service for authentification
// It implements the AuthService interface
type AuthServiceImpl struct {
	ctx    context.Context
	gormDb *gorm.DB
}

// AuthServiceImpl implements the AuthService interface
var _ AuthService = &AuthServiceImpl{}

// NewAuthService creates a new service used for authentification
func NewAuthService(ctx context.Context, gormDb *gorm.DB) AuthService {
	return &AuthServiceImpl{ctx: ctx, gormDb: gormDb}
}

// UserSignUp creates a new user
func (s *AuthServiceImpl) UserSignUp(user *models.UserSignUp) (*models.User, error) {

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Create a new user from the UserSignUp input
	newUser := &models.User{
		//Name:     user.Name,
		Email:    strings.ToLower(user.Email),
		Password: hashedPassword,
		Role:     models.RoleUser,
		Verified: false,
	}

	// Add the new user to the database
	if results := s.gormDb.Create(&newUser); results.Error != nil {
		var pgError *pgconn.PgError
		if errors.As(results.Error, &pgError) && errors.Is(results.Error, pgError) && pgError.Code == "23505" {
			return nil, errors.New("user with that email already exist")
		}
		return nil, results.Error
	}

	return newUser, nil
}
