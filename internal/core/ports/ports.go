package ports

import "github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"

type UserService interface {
	Find(id string) (models.User, error)
	Save(u models.User) (models.User, error)
}

type UserRepo interface {
	FindById(id string) (models.User, error)
	FindAll() ([]models.User, error)
	Save(u models.User) (models.User, error)
}
