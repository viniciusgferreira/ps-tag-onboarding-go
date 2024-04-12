package service

import (
	"context"
	"errors"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"log/slog"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("user with the same first and last name already exists")
)

type UserRepository interface {
	FindById(ctx context.Context, id string) (*model.User, error)
	Save(ctx context.Context, u model.User) (*model.User, error)
	Update(ctx context.Context, u model.User) (*model.User, error)
	ExistsByFirstNameAndLastName(ctx context.Context, u model.User) (bool, error)
}
type Service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

type ValidationError struct {
	Message string
	Details []string
}

func (r ValidationError) Error() string {
	return r.Message
}

func (s *Service) FindById(ctx context.Context, id string) (*model.User, error) {
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

func (s *Service) Save(ctx context.Context, u model.User) (*model.User, error) {
	usernameTaken, err := s.repo.ExistsByFirstNameAndLastName(ctx, u)
	if err != nil {
		return nil, err
	}
	if usernameTaken {
		return nil, ErrUsernameTaken
	}
	savedUser, err := s.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

func (s *Service) Update(ctx context.Context, updatedUser model.User) (*model.User, error) {
	existingUser, err := s.repo.FindById(ctx, updatedUser.ID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}
	usernameTaken, err := s.repo.ExistsByFirstNameAndLastName(ctx, updatedUser)
	if err != nil {
		return nil, err
	}
	if usernameTaken {
		return nil, ErrUsernameTaken
	}
	updatedUserResult, err := s.repo.Update(ctx, updatedUser)
	if err != nil {
		return nil, err
	}
	return updatedUserResult, nil
}
