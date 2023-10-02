package services

import (
	"context"
	"errors"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"strings"
)

// IUserService is an interface for the UserService
// It declares the methods that the UserService must implement
type IUserService interface {
	Create(user *models.UserSignUp) (*models.User, error)
	FindById(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}

// UserService is the service for the user
// It implements the IUserService interface
type UserService struct {
	ctx    context.Context
	gormDb *gorm.DB
}

// UserService implements the IUserService interface
var _ IUserService = &UserService{}

func NewUserService(ctx context.Context, gormDb *gorm.DB) IUserService {
	return &UserService{ctx: ctx, gormDb: gormDb}
}

func (s *UserService) Create(user *models.UserSignUp) (*models.User, error) {

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

func (s *UserService) FindById(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := s.gormDb.Find(&user, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &user, nil
	}
	return nil, nil
}

func (s *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.gormDb.Find(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &user, nil
	}
	return nil, nil
}
