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

// IAuthService is the interface that the service must implement
// It declares the methods that the service must implement
type IAuthService interface {
	UserSignUp(user *models.UserSignUp) (*models.User, error)
}

// AuthService is the service for authentification
// It implements the IAuthService interface
type AuthService struct {
	ctx    context.Context
	gormDb *gorm.DB
}

// AuthService implements the IAuthService interface
var _ IAuthService = &AuthService{}

func NewAuthService(ctx context.Context, gormDb *gorm.DB) IAuthService {
	return &AuthService{ctx: ctx, gormDb: gormDb}
}

func (s *AuthService) UserSignUp(user *models.UserSignUp) (*models.User, error) {

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Create a new user from the UserSignUp input
	newUser := &models.User{
		Name:     user.Name,
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
