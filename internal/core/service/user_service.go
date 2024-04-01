package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"log/slog"
	"regexp"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with the same first and last name already exists")
	ErrInvalidAge        = errors.New("user must be at least 18 years old")
	ErrInvalidName       = errors.New("user first and last name cannot be empty")
	ErrInvalidEmail      = errors.New("invalid email")
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

func NewUserService(repo UserRepository) *UserService {
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
		Message: "user did not pass validation",
		Details: details,
	}
}

func (r ValidationError) Error() string {
	return r.Message
}

func (s *UserService) Find(ctx *gin.Context, id string) (*model.User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		slog.Warn("user not found")
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Save(ctx *gin.Context, u model.User) (*model.User, error) {
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		slog.Error("validationError", "error", validationErrors)
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
		slog.Error("userService.Save", "error", err.Error())
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx *gin.Context, u model.User) (*model.User, error) {
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		slog.Error("validationError", "error", validationErrors)
		return nil, NewValidationErrorWithDetails(validationErrors)
	}
	user, err := s.repo.FindById(ctx, u.ID)
	if err != nil {
		return nil, err
	}
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
		slog.Error("userService.Update", "error", err.Error())
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

func (s *UserService) validateName(firstName, lastName string) error {
	if len(firstName) <= 0 || len(lastName) <= 0 {
		return ErrInvalidName
	}
	return nil
}

func (s *UserService) validateEmail(email string) error {
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
