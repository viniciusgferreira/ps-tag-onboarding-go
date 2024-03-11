package httpserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Message string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

type UserHandler struct {
	service service.UserService
}

func New(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/users/:id", Handler: h.GetByID},
		{Method: http.MethodPost, Path: "/users", Handler: h.Create},
		{Method: http.MethodPut, Path: "/users/:id", Handler: h.Update},
	}
}

// GetByID godoc
// @Summary Find user by ID
// @Description Get user based on request path
// @Tags users
// @Produce json
// @Param id path string true "id"
// @Sucess 201 {object} model.User
// @Failure 404
// @Failure 400
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(ctx *gin.Context) {
	user, err := h.service.Find(ctx, ctx.Param("id"))
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
// @Failure 409
// @Failure 400
// @Router /users [post]
func (h *UserHandler) Create(ctx *gin.Context) {
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	savedUser, err := h.service.Save(ctx, user)
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
// @Failure 404 {object} model.User "User not found"
// @Failure 400 {object} model.User "Invalid id"
// @Failure 409 "User already exists"
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *gin.Context) {
	u := model.User{}
	u.ID = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	updatedUser, err := h.service.Update(ctx, u)
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
	case errors.Is(err, service.ErrUserAlreadyExists):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	default:
		slog.Error("handler", "checkErr", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal Server Error"})
	}
}
