package httpserver

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := &UserMockService{}
	handler := NewUserHandler(mockUserService)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	// assert all expectations, reset recorder and create new context
	reset := func() {
		mockUserService.AssertExpectations(t)
		mockUserService.ExpectedCalls = nil
		recorder = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(recorder)
	}

	t.Run("Save user sucessfully", func(t *testing.T) {
		t.Cleanup(reset)
		user := model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Age:       18,
		}
		jsonBody, _ := json.Marshal(user)
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonBody)),
		}
		mockUserService.On("Save", ctx, user).Return(&user, nil).Once()

		handler.Create(ctx)
		var responseUser model.User
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, http.StatusCreated, ctx.Writer.Status())
		assert.Equal(t, user.FirstName, responseUser.FirstName)
		assert.Equal(t, user.LastName, responseUser.LastName)
		assert.Equal(t, user.Email, responseUser.Email)
		assert.Equal(t, user.Age, responseUser.Age)
	})

	t.Run("User already exists Error", func(t *testing.T) {
		t.Cleanup(reset)
		user := model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@email.com",
			Age:       18,
		}
		jsonBody, _ := json.Marshal(user)
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonBody)),
		}
		err := service.ErrUserAlreadyExists
		mockUserService.On("Save", ctx, user).Return(nil, err).Once()

		handler.Create(ctx)
		var responseBody ErrorResponse
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseBody)

		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
		assert.NotEmpty(t, responseBody.Message)
		assert.Equal(t, err.Error(), responseBody.Message)
		assert.Nil(t, responseBody.Details)
	})
}

func TestUserHandler_Find(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := &UserMockService{}
	handler := NewUserHandler(mockUserService)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	// assert all expectations, reset recorder and create new context
	reset := func() {
		mockUserService.AssertExpectations(t)
		mockUserService.ExpectedCalls = nil
		recorder = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(recorder)
	}
	t.Run("Find User sucessfully", func(t *testing.T) {
		t.Cleanup(reset)
		id := primitive.NewObjectID().String()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: id,
		}}
		user := model.User{
			ID:        id,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@email.com",
			Age:       18,
		}
		mockUserService.On("FindById", ctx, id).Return(&user, nil).Once()

		handler.FindById(ctx)
		var responseUser model.User
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, user.ID, responseUser.ID)
		assert.Equal(t, user.FirstName, responseUser.FirstName)
		assert.Equal(t, user.LastName, responseUser.LastName)
		assert.Equal(t, user.Email, responseUser.Email)
		assert.Equal(t, user.Age, responseUser.Age)
	})
	t.Run("User not found", func(t *testing.T) {
		t.Cleanup(reset)
		id := primitive.NewObjectID().String()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: id,
		}}
		err := service.ErrUserNotFound
		mockUserService.On("FindById", ctx, id).Return(nil, err).Once()

		handler.FindById(ctx)
		var responseBody ErrorResponse
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseBody)

		assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
		assert.NotEmpty(t, responseBody.Message)
		assert.Equal(t, err.Error(), responseBody.Message)
		assert.Nil(t, responseBody.Details)
	})
}

func TestUserHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := &UserMockService{}
	handler := NewUserHandler(mockUserService)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	// assert all expectations, reset recorder and create new context
	reset := func() {
		mockUserService.AssertExpectations(t)
		mockUserService.ExpectedCalls = nil
		recorder = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(recorder)
	}
	t.Run("Update User sucessfully", func(t *testing.T) {
		t.Cleanup(reset)
		id := primitive.NewObjectID().String()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: id,
		}}
		updatedUser := model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Age:       18,
		}
		jsonBody, _ := json.Marshal(updatedUser)
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonBody)),
		}
		updatedUser.ID = id
		mockUserService.On("Update", ctx, updatedUser).Return(&updatedUser, nil).Once()

		handler.Update(ctx)
		var responseUser model.User
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, updatedUser.ID, responseUser.ID)
		assert.Equal(t, updatedUser.FirstName, responseUser.FirstName)
		assert.Equal(t, updatedUser.LastName, responseUser.LastName)
		assert.Equal(t, updatedUser.Email, responseUser.Email)
		assert.Equal(t, updatedUser.Age, responseUser.Age)
	})
	t.Run("User not found", func(t *testing.T) {
		t.Cleanup(reset)
		id := primitive.NewObjectID().String()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: id,
		}}
		updatedUser := model.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Age:       18,
		}
		jsonBody, _ := json.Marshal(updatedUser)
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonBody)),
		}
		updatedUser.ID = id
		err := service.ErrUserNotFound
		mockUserService.On("Update", ctx, updatedUser).Return(nil, err).Once()

		handler.Update(ctx)
		var responseBody ErrorResponse
		_ = json.Unmarshal(recorder.Body.Bytes(), &responseBody)

		assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
		assert.NotEmpty(t, responseBody.Message)
		assert.Equal(t, err.Error(), responseBody.Message)
		assert.Nil(t, responseBody.Details)
	})
}
