package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	repo ports.UserRepository
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with same first and last name already exists")
)

func New(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) Find(ctx *gin.Context, id string) (*models.User, error) {
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
	exists, err := s.repo.ExistsByFirstNameAndLastName(ctx, u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}
	u.ID = uuid.NewString()
	user, err := s.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx *gin.Context, u models.User) (*models.User, error) {
	_, err := s.repo.FindById(ctx, u.ID)
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
