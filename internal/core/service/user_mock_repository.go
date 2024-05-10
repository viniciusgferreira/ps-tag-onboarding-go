package service

import (
	"context"
	"errors"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
)

type MockUserRepository struct {
	Users []model.User
}

func (m *MockUserRepository) FindById(_ context.Context, id string) (*model.User, error) {
	for _, user := range m.Users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
func (m *MockUserRepository) Save(_ context.Context, u model.User) (*model.User, error) {
	m.Users = append(m.Users, u)
	return &u, nil
}
func (m *MockUserRepository) Update(_ context.Context, updatedUser model.User) (*model.User, error) {
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

func (m *MockUserRepository) ExistsByFirstNameAndLastName(_ context.Context, u model.User) (bool, error) {
	for _, user := range m.Users {
		if user.ID != u.ID && user.FirstName == u.FirstName && user.LastName == u.LastName {
			return true, nil
		}
	}
	return false, nil
}
