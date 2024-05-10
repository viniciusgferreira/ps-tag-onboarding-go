package httpserver

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
)

type UserMockService struct {
	mock.Mock
}

func (m *UserMockService) Save(ctx context.Context, user model.User) (*model.User, error) {
	called := m.Called(ctx, user)
	if len(called) == 0 {
		panic("no return value specified for Save")
	}
	resultUser := called.Get(0)
	err := called.Error(1)
	if resultUser != nil {
		return resultUser.(*model.User), err
	}
	return nil, err
}
func (m *UserMockService) Update(ctx context.Context, user model.User) (*model.User, error) {
	called := m.Called(ctx, user)
	if len(called) == 0 {
		panic("no return value specified for Update")
	}
	resultUser := called.Get(0)
	err := called.Error(1)
	if resultUser != nil {
		return resultUser.(*model.User), err
	}
	return nil, err
}
func (m *UserMockService) FindById(ctx context.Context, id string) (*model.User, error) {
	called := m.Called(ctx, id)
	if len(called) == 0 {
		panic("no return value specified for FindById")
	}
	resultUser := called.Get(0)
	err := called.Error(1)
	if resultUser != nil {
		return resultUser.(*model.User), err
	}
	return nil, err
}
