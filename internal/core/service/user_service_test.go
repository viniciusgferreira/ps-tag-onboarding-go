package service

import (
	"errors"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func TestUserService(t *testing.T) {
	user := model.User{
		ID:        primitive.NewObjectID().Hex(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Age:       21,
	}
	t.Run("Save user", func(t *testing.T) {
		mockRepo := &test.MockUserRepository{}
		service := New(mockRepo)

		result, err := service.Save(nil, user)
		if err != nil {
			t.Errorf("error saving user: %v", err)
		}
		if !reflect.DeepEqual(result, &user) {
			t.Errorf("expected: %v, result: %v", user, result)
		}
		if len(mockRepo.Users) != 1 || !reflect.DeepEqual(&mockRepo.Users[0], result) {
			t.Errorf("user not saved properly in mock repository")
		}
	})

	t.Run("Find user", func(t *testing.T) {
		Users := []model.User{user}
		mockRepo := &test.MockUserRepository{Users: Users}
		service := New(mockRepo)

		result, err := service.Find(nil, user.ID)
		if err != nil {
			t.Errorf("error finding user: %v", err)
		}
		if !reflect.DeepEqual(result, &user) {
			t.Errorf("expected: %v, result: %v", user, result)
		}
	})
	t.Run("Update user", func(t *testing.T) {
		Users := []model.User{user}
		updatedUser := model.User{
			ID:        user.ID,
			FirstName: "Doe",
			LastName:  "New",
			Email:     "new@doe.com",
			Age:       23,
		}
		mockRepo := &test.MockUserRepository{Users: Users}
		service := &UserService{repo: mockRepo}

		result, err := service.Update(nil, updatedUser)
		if err != nil {
			t.Errorf("error finding user: %v", err)
		}
		if !reflect.DeepEqual(result, &updatedUser) {
			t.Errorf("expected: %v, result: %v", updatedUser, result)
		}
		if len(mockRepo.Users) != 1 {
			t.Errorf("user not updated properly in mock repository")
		}
	})
}
func TestUserValidation(t *testing.T) {
	t.Run("Validate user", func(t *testing.T) {
		user := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Age:       21,
		}
		mockRepo := &test.MockUserRepository{}
		service := New(mockRepo)

		err := service.Validate(user)
		if err != nil {
			t.Errorf("error validating user: %v", err)
		}
	})
	t.Run("Should return error with invalid user", func(t *testing.T) {
		invalidUser := model.User{}
		mockRepo := &test.MockUserRepository{}
		service := New(mockRepo)

		err := service.Validate(invalidUser)
		if err == nil {
			t.Errorf("validation should fail: %v", err)
		}
	})
	t.Run("Should return error with invalid age", func(t *testing.T) {
		userWithInvalidAge := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Age:       17,
		}
		mockRepo := &test.MockUserRepository{}
		service := New(mockRepo)

		errs := service.Validate(userWithInvalidAge)
		if len(errs) == 0 || !errors.Is(errs[0], ErrInvalidAge) {
			t.Errorf("expected err: %v, got: %v\n", ErrInvalidAge, errs[0])
		}
	})
	t.Run("Should return error with invalid email", func(t *testing.T) {
		userWithInvalidEmail := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe2.com",
			Age:       20,
		}
		mockRepo := &test.MockUserRepository{}
		service := New(mockRepo)

		errs := service.Validate(userWithInvalidEmail)
		if len(errs) != 1 {
			t.Errorf("expected err: %v, got: %v\n", ErrInvalidEmail, nil)
			return
		}
		if !errors.Is(ErrInvalidEmail, errs[0]) {
			t.Errorf("expected err: %v, got: %v\n", ErrInvalidEmail, errs[0])
		}
	})
}
