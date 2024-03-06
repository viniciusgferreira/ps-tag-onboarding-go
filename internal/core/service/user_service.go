package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User with same first and last name already exists")
	ErrInvalidAge        = errors.New("User must be at least 18 years old")
	ErrInvalidName       = errors.New("User first and last name cannot be empty")
	ErrInvalidEmail      = errors.New("Invalid email format")
	ErrInvalidID         = errors.New("Invalid id (ObjectID)")
)

type UserService struct {
	repo ports.UserRepository
}

func New(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type ValidationError struct {
	Message string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

func NewValidationError(validationErrors []error) ValidationError {
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

func (s *UserService) Find(ctx *gin.Context, id string) (*models.User, error) {
	err := s.validateID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Save(ctx *gin.Context, u models.User) (*models.User, error) {
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		return nil, NewValidationError(validationErrors)
	}
	exists, err := s.repo.ExistsByFirstNameAndLastName(ctx, u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err := s.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx *gin.Context, u models.User) (*models.User, error) {
	err := s.validateID(u.ID)
	if err != nil {
		return nil, err
	}
	validationErrors := s.Validate(u)
	if validationErrors != nil {
		return nil, NewValidationError(validationErrors)
	}
	_, err = s.repo.FindById(ctx, u.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	exists, err := s.repo.ExistsByFirstNameAndLastName(ctx, u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}
	user, err := s.repo.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) Validate(u models.User) []error {
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
		return ErrInvalidID
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
