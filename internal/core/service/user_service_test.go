package service

import (
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
	"testing"
)

var (
	user, _   = model.NewUser(primitive.NewObjectID().Hex(), "John", "Doe", "john@doe.com", 21)
	validUser = *user
)

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
	t.Run("Should return error for user that already exists", func(t *testing.T) {
		Users := []model.User{validUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := NewUserService(mockRepo)

		newUser, _ := model.NewUser(primitive.NewObjectID().Hex(), "John", "Doe", "doe@j.com", 20)

		_, err := service.Save(nil, *newUser)
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

		result, err := service.FindById(nil, validUser.ID)
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

		result, err := service.FindById(nil, validUser.ID)
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
		service := &Service{repo: mockRepo}

		updatedUser, _ := model.NewUser(validUser.ID, "Doe", "NewUserService", "new@doe.com", 23)

		result, err := service.Update(nil, *updatedUser)
		if err != nil {
			t.Errorf("error finding user: %v", err)
		}
		if !reflect.DeepEqual(result, updatedUser) {
			t.Errorf("expected: %v, result: %v", updatedUser, result)
		}
		if len(mockRepo.Users) != 1 {
			t.Errorf("user not updated properly in mock repository")
		}
	})
	t.Run("return error on update with first and lastname that already exists", func(t *testing.T) {
		updatedUser, _ := model.NewUser(primitive.NewObjectID().Hex(), "John", "Doe Second", "new@doe.com", 25)
		Users := []model.User{validUser, *updatedUser}
		mockRepo := &MockUserRepository{Users: Users}
		service := &Service{repo: mockRepo}

		changedValidUser, _ := model.NewUser(validUser.ID, updatedUser.FirstName, updatedUser.LastName, "new@doe.com", 25)

		_, err := service.Update(nil, *changedValidUser)
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
		service := &Service{repo: mockRepo}

		updatedUser, _ := model.NewUser(primitive.NewObjectID().Hex(), "John", "Doe", "new@doe.com", 23)

		result, err := service.Update(nil, *updatedUser)
		if err == nil || result != nil {
			t.Errorf("update should fail: %v", err)
		}
	})
}
