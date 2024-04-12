package httpserver

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/dto"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
	"net/http"
)

var validator = galidator.New().Validator(dto.UserInput{})

type ErrorResponse struct {
	Message string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}

type UserService interface {
	Save(ctx context.Context, u model.User) (*model.User, error)
	Update(ctx context.Context, u model.User) (*model.User, error)
	FindById(ctx context.Context, id string) (*model.User, error)
}

func (h *UserHandler) SetupRoutes(r *Router) {
	r.Handle(http.MethodGet, "/users/:id", h.FindById)
	r.Handle(http.MethodPost, "/users", h.Create)
	r.Handle(http.MethodPut, "/users/:id", h.Update)
}

// FindById godoc
// @Summary Find user by ID
// @Description Get user based on request path
// @Tags users
// @Produce json
// @Param id path string true "id"
// @Sucess 201 {object} model.User
// @Failure 404
// @Failure 400
// @Router /users/{id} [get]
func (h *UserHandler) FindById(ctx *gin.Context) {
	user, err := h.service.FindById(ctx, ctx.Param("id"))
	if err != nil {
		checkErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// Create godoc
// @Summary Create a new user
// @Description Create new user based on request body input
// @Tags users
// @Accept json
// @Produce json
// @Param User body model.User true "User input"
// @Sucess 201 {object} modes.User
// @Failure 400
// @Router /users [post]
func (h *UserHandler) Create(ctx *gin.Context) {
	userInput := dto.UserInput{}
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "user did not pass validation",
			Details: validator.DecryptErrors(err),
		})
		return
	}
	user, err := model.NewUser(ctx.Param("id"), userInput.FirstName, userInput.LastName, userInput.Email, userInput.Age)
	savedUser, err := h.service.Save(ctx, *user)
	if err != nil {
		checkErr(ctx, err)
		return
	}
	ctx.IndentedJSON(http.StatusCreated, savedUser)
}

// Update godoc
// @Summary Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param User body model.User true "User input"
// @Sucess 200 {object} model.User
// @Failure 404 {object} model.User
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *gin.Context) {
	userInput := dto.UserInput{}
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "user did not pass validation",
			Details: validator.DecryptErrors(err),
		})
		return
	}
	user, err := model.NewUser(ctx.Param("id"), userInput.FirstName, userInput.LastName, userInput.Email, userInput.Age)
	updatedUser, err := h.service.Update(ctx, *user)
	if err != nil {
		checkErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

func checkErr(ctx *gin.Context, err error) {
	var validationErr service.ValidationError
	switch {
	case errors.As(err, &validationErr):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: validationErr.Message, Details: validationErr.Details})
	case errors.Is(err, service.ErrUserNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	case errors.Is(err, service.ErrUsernameTaken):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	default:
		slog.Error(ctx.Request.RequestURI, "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal Server Error"})
	}
}
