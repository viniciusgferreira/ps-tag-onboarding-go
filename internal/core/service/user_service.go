package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/port"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User with the same first and last name already exists")
	ErrInvalidAge        = errors.New("User must be at least 18 years old")
	ErrInvalidName       = errors.New("User first and last name cannot be empty")
	ErrInvalidEmail      = errors.New("Invalid email format")
)

type UserService struct {
	repo port.UserRepository
}

func New(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type ValidationError struct {
	Message string
	Details []string
}

func NewValidationErrorWithDetails(validationErrors []error) ValidationError {
	details := make([]string, len(validationErrors))
	for i, err := range validationErrors {
		details[i] = err.Error()
	}
	return ValidationError{
		Message: "User did not pass validation",
		Details: details,
	}
}

func (r ValidationError) Error() string {
	return r.Message
}

func (s *UserService) Find(ctx *gin.Context, id string) (*model.User, error) {
	err := s.validateID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindById(ctx, id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Save(ctx *gin.Context, u model.User) (*model.User, error) {
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		return nil, NewValidationErrorWithDetails(validationErrors)
	}
	exists := s.repo.ExistsByFirstNameAndLastName(ctx, u)
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err := s.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx *gin.Context, u model.User) (*model.User, error) {
	err := s.validateID(u.ID)
	if err != nil {
		return nil, err
	}
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		return nil, NewValidationErrorWithDetails(validationErrors)
	}
	user, _ := s.repo.FindById(ctx, u.ID)
	if user == nil {
		return nil, ErrUserNotFound
	}
	exists := s.repo.ExistsByFirstNameAndLastName(ctx, u)
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err = s.repo.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) Validate(u model.User) []error {
	var validationErrors []error
	if err := s.validateName(u.FirstName, u.LastName); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateEmail(u.Email); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateAge(u.Age); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
func (s *UserService) validateID(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserService) validateName(firstName, lastName string) error {
	if len(firstName) <= 0 || len(lastName) <= 0 {
		return ErrInvalidName
	}
	return nil
}

func (s *UserService) validateEmail(email string) error {
	if len(email) <= 0 {
		return errors.New("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func (s *UserService) validateAge(age int) error {
	if age < 18 {
		return ErrInvalidAge
	}
	return nil
}
