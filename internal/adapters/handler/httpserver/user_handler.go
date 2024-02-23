package httpserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"net/http"
)

type UserHandler struct {
	service ports.UserService
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
// @Sucess 201 {object} models.User
// @Failure 404
// @Failure 400
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(ctx *gin.Context) {
	user, err := h.service.Find(ctx, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
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
// @Param User body models.User true "User input"
// @Sucess 201 {object} modes.User
// @Failure 409
// @Failure 400
// @Router /users [post]
func (h *UserHandler) Create(ctx *gin.Context) {
	user := models.User{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := user.Validate()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	savedUser, err := h.service.Save(ctx, user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
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
// @Param User body models.User true "User input"
// @Sucess 200 {object} models.User
// @Failure 404 {object} models.User "User not found"
// @Failure 400 {object} models.User "Invalid id"
// @Failure 409 "User already exists"
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *gin.Context) {
	u := models.User{}
	if err := ctx.BindJSON(&u); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	u.ID = ctx.Param("id")
	err := u.Validate()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	updatedUser, err := h.service.Update(ctx, u)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUserAlreadyExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

func New(s ports.UserService) *UserHandler {
	return &UserHandler{service: s}
}
