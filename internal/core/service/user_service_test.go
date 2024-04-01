package service

import (
	"errors"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
	"testing"
)

var validUser = model.User{
	ID:        primitive.NewObjectID().Hex(),
	FirstName: "John",
	LastName:  "Doe",
	Email:     "john@doe.com",
	Age:       21,
}

func TestSaveUser(t *testing.T) {
	t.Run("Save valid user", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		result, err := service.Save(nil, validUser)
		if err != nil {
			t.Errorf("error saving user: %v", err)
		}
		if !reflect.DeepEqual(result, &validUser) {
			t.Errorf("expected: %v, result: %v", validUser, result)
		}
		if len(mockRepo.Users) != 1 || !reflect.DeepEqual(&mockRepo.Users[0], result) {
			t.Errorf("user not saved properly in mock repository")
		}
	})
	t.Run("Should return error with empty user", func(t *testing.T) {
		emptyUser := model.User{}
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		result, err := service.Save(nil, emptyUser)
		if err == nil || result != nil {
			t.Errorf("error saving user: %v", err)
		}
		if len(mockRepo.Users) != 0 {
			t.Errorf("empty user should not be saved in mock repository")
		}
	})
	t.Run("Should return error for user that already exists", func(t *testing.T) {
		Users := []model.User{validUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := NewUserService(mockRepo)

		newUser := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "doe@j.com",
			Age:       20,
		}

		_, err := service.Save(nil, newUser)
		if err == nil {
			t.Errorf("should return error saving user: %v", err)
		}
		if len(mockRepo.Users) != 1 {
			log.Println(mockRepo.Users)
			t.Errorf("duplicate user should not be saved in mock repository")
		}
	})
}
func TestFindUser(t *testing.T) {
	t.Run("Find user", func(t *testing.T) {
		Users := []model.User{validUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := NewUserService(mockRepo)

		result, err := service.Find(nil, validUser.ID)
		if err != nil {
			t.Errorf("error finding user: %v", err)
		}
		if !reflect.DeepEqual(result, &validUser) {
			t.Errorf("expected: %v, result: %v", validUser, result)
		}
	})
	t.Run("Should return error with user not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		result, err := service.Find(nil, validUser.ID)
		if result != nil {
			t.Errorf("should return userNotFound")
		}
		if err == nil {
			t.Errorf("expected: %v, result: %v", ErrUserNotFound, err)
		}
	})
}
func TestUpdateUser(t *testing.T) {
	t.Run("Update valid user", func(t *testing.T) {
		Users := []model.User{validUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := &UserService{repo: mockRepo}

		updatedUser := model.User{
			ID:        validUser.ID,
			FirstName: "Doe",
			LastName:  "NewUserService",
			Email:     "new@doe.com",
			Age:       23,
		}

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
	t.Run("return error on update with first and lastname that already exists", func(t *testing.T) {
		updatedUser := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe Second",
			Email:     "new@doe.com",
			Age:       25,
		}
		Users := []model.User{validUser, updatedUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := &UserService{repo: mockRepo}

		changedValidUser := model.User{
			ID:        validUser.ID,
			FirstName: updatedUser.FirstName, // already exists
			LastName:  updatedUser.LastName,  // already exists
			Email:     "new@doe.com",
			Age:       25,
		}

		_, err := service.Update(nil, changedValidUser)
		if err == nil {
			t.Errorf("should return error updating user: %v", err)
		}
		if len(mockRepo.Users) != 2 {
			log.Println(mockRepo.Users)
			t.Errorf("duplicate user should not be updated in mock repository")
		}
	})
	t.Run("return error with existing user first and last name", func(t *testing.T) {
		Users := []model.User{validUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := &UserService{repo: mockRepo}

		updatedUser := model.User{
			ID:        primitive.NewObjectID().Hex(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "new@doe.com",
			Age:       23,
		}

		result, err := service.Update(nil, updatedUser)
		if err == nil || result != nil {
			t.Errorf("update should fail: %v", err)
		}
	})
}
func TestUserValidation(t *testing.T) {
	t.Run("Validate user without error", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		err := service.Validate(validUser)
		if err != nil {
			t.Errorf("error validating user: %v", err)
		}
	})
	t.Run("Should return error with invalid user", func(t *testing.T) {
		invalidUser := model.User{}
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		err := service.Validate(invalidUser)
		if err == nil || reflect.TypeOf(err).String() != "[]error" {
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
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

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
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

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
