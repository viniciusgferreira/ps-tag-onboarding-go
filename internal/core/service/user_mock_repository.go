package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
)

type MockUserRepository struct {
	Users []model.User
}

func (m *MockUserRepository) FindById(ctx *gin.Context, id string) (*model.User, error) {
	for _, user := range m.Users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
func (m *MockUserRepository) Save(ctx *gin.Context, u model.User) (*model.User, error) {
	m.Users = append(m.Users, u)
	return &u, nil
}
func (m *MockUserRepository) Update(ctx *gin.Context, updatedUser model.User) (*model.User, error) {
	index := -1
	for i, user := range m.Users {
		if user.ID == updatedUser.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errors.New("user not found")
	}
	m.Users[index] = updatedUser
	return &m.Users[index], nil
}

func (m *MockUserRepository) ExistsByFirstNameAndLastName(ctx *gin.Context, u model.User) bool {
	return false
}
