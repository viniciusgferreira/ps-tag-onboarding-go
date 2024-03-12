package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"regexp"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User with the same first and last name already exists")
	ErrInvalidAge        = errors.New("User must be at least 18 years old")
	ErrInvalidName       = errors.New("User first and last name cannot be empty")
	ErrInvalidEmail      = errors.New("Invalid email format")
)

type UserRepository interface {
	FindById(ctx *gin.Context, id string) (*model.User, error)
	Save(ctx *gin.Context, u model.User) (*model.User, error)
	Update(ctx *gin.Context, u model.User) (*model.User, error)
	ExistsByFirstNameAndLastName(ctx *gin.Context, u model.User) (bool, error)
}
type UserService struct {
	repo UserRepository
}

func New(repo UserRepository) *UserService {
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
	slog.Info("UserService.Find", "id", id)
	err := s.validateID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindById(ctx, id)
	if err != nil || user == nil {
		slog.Warn("User not found")
		return nil, ErrUserNotFound
	}
	slog.Info("UserService.Find completed sucessfully")
	return user, nil
}

func (s *UserService) Save(ctx *gin.Context, u model.User) (*model.User, error) {
	slog.Info("UserService.Save - Saving new user")
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		return nil, NewValidationErrorWithDetails(validationErrors)
	}
	exists, err := s.repo.ExistsByFirstNameAndLastName(ctx, u)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err := s.repo.Save(ctx, u)
	if err != nil {
		slog.Error("UserService.Save failed")
		return nil, err
	}
	slog.Info("UserService.Save completed sucessfully")
	return user, nil
}

func (s *UserService) Update(ctx *gin.Context, u model.User) (*model.User, error) {
	slog.Info("UserService.Update", "Updating user", u.ID)
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
	exists, err := s.repo.ExistsByFirstNameAndLastName(ctx, u)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err = s.repo.Update(ctx, u)
	if err != nil {
		slog.Error("UserService.Update failed")
		return nil, err
	}
	slog.Info("UserService.Update completed sucessfully")
	return user, nil
}
func (s *UserService) Validate(u model.User) []error {
	var validationErrors []error
	if err := s.validateName(u.FirstName, u.LastName); err != nil {
		slog.Warn("UserService.Validate", "name", "Invalid name")
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateEmail(u.Email); err != nil {
		slog.Warn("UserService.Validate", "email", "Invalid email")
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateAge(u.Age); err != nil {
		slog.Warn("UserService.Validate", "age", "Invalid age")
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		slog.Error("UserService.Validate", "Validation errors", len(validationErrors))
		return validationErrors
	}
	slog.Info("User is valid")
	return nil
}
func (s *UserService) validateID(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Mongodb", "ObjectID conversion", err.Error())
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
