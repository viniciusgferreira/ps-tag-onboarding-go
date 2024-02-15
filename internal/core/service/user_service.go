package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
)

type Service struct {
	repo ports.UserRepo
}

func New(repo ports.UserRepo) *Service {
	return &Service{repo: repo}
}
func (s *Service) Find(ctx *gin.Context, id string) (*models.User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service) Save(ctx *gin.Context, u *models.User) (*models.User, error) {
	u.ID = uuid.NewString()
	user, err := s.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}