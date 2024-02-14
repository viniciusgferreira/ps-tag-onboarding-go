package service

import (
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
)

type Service struct {
	repo ports.UserRepo
}

func New(repo ports.UserRepo) *Service {
	return &Service{repo: repo}
}
func (s *Service) Find(id string) (models.User, error) {
	return s.repo.FindById(id)
}

func (s *Service) Save(u models.User) (models.User, error) {
	return models.User{}, nil
}
