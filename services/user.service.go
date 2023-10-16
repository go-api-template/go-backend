package services

import (
	"context"
	"errors"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/utils"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"strings"
)

// UserService is an interface for the UserServiceImpl
// It declares the methods that the UserServiceImpl must implement
type UserService interface {
	Create(user *models.UserSignUp) (*models.User, error)

	FindAll(params httputil.Filter) ([]models.User, error)
	FindById(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByVerificationToken(verificationToken string) (*models.User, error)
	FindByResetPasswordToken(resetPasswordToken string) (*models.User, error)

	Update(id uuid.UUID, user *models.User) (*models.User, error)
	Delete(id uuid.UUID) error
}

// UserServiceImpl is the service for the user
// It implements the UserService interface
type UserServiceImpl struct {
	ctx    context.Context
	gormDb *gorm.DB
}

// UserServiceImpl implements the UserService interface
var _ UserService = &UserServiceImpl{}

func NewUserService(ctx context.Context, gormDb *gorm.DB) UserService {
	return &UserServiceImpl{ctx: ctx, gormDb: gormDb}
}

func (s *UserServiceImpl) Create(user *models.UserSignUp) (*models.User, error) {

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Verification code
	verificationCode := utils.GenerateRandomString(32)
	verificationCode = utils.Encode(verificationCode)

	// Create a new user from the SignUp input
	newUser := &models.User{
		//Name:             user.Name,
		Email:             strings.ToLower(user.Email),
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Password:          hashedPassword,
		Role:              models.RoleUser,
		VerificationToken: verificationCode,
		Verified:          false,
	}

	// Set the name if it is empty
	if len(newUser.Name) == 0 {
		newUser.Name = newUser.FirstName + " " + newUser.LastName
	}

	// If this is the first user, set the role to admin
	if s.gormDb.First(&models.User{}).RowsAffected == 0 {
		newUser.Role = models.RoleAdmin
		newUser.Verified = true
		newUser.VerificationToken = ""
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

func (s *UserServiceImpl) FindAll(params httputil.Filter) ([]models.User, error) {
	var users []models.User
	result := s.gormDb.Scopes(params.Apply).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return users, nil
	}
	return nil, nil
}

func (s *UserServiceImpl) FindById(id uuid.UUID) (*models.User, error) {
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

func (s *UserServiceImpl) FindByEmail(email string) (*models.User, error) {
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

func (s *UserServiceImpl) FindByVerificationToken(verificationCode string) (*models.User, error) {
	var user models.User
	result := s.gormDb.Find(&user, "verification_token = ?", verificationCode)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &user, nil
	}
	return nil, nil
}

func (s *UserServiceImpl) FindByResetPasswordToken(resetPasswordToken string) (*models.User, error) {
	var user models.User
	result := s.gormDb.Find(&user, "reset_token = ?", resetPasswordToken)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &user, nil
	}
	return nil, nil
}

func (s *UserServiceImpl) Update(id uuid.UUID, user *models.User) (*models.User, error) {
	// Set the verification code if the user is not verified yet and the verification code is empty
	if !user.Verified && user.VerificationToken == "" {
		user.VerificationToken = utils.Encode(utils.GenerateRandomString(32))
	}
	// Remove the verification code if the user is verified and the verification code is not empty
	if user.Verified && user.VerificationToken != "" {
		user.VerificationToken = ""
	}

	// Update the user
	// Use select("*") to update all fields, even if they are empty
	// this prevent zero value fields from being updated
	result := s.gormDb.Model(user).Select("*").Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return s.FindById(id)
	}
	return nil, nil
}

func (s *UserServiceImpl) Delete(id uuid.UUID) error {
	// Anonymize the user
	user, err := s.anonymize(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("unknown user")
	}

	// Delete the user from the database
	result := s.gormDb.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}
	return errors.New("unknown user")
}

func (s *UserServiceImpl) anonymize(id uuid.UUID) (*models.User, error) {
	// Get the user from the database
	user, err := s.FindById(id)
	if err != nil {
		return nil, err
	}

	// Anonymize the user
	user.Email = faker.Email()
	user.Name = faker.Name()
	user.FirstName = faker.FirstName()
	user.LastName = faker.LastName()

	return s.Update(id, user)
}
