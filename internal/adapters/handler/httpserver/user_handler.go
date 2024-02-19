package httpserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository/mongo"
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

func (h *UserHandler) GetByID(ctx *gin.Context) {
	user, err := h.service.Find(ctx, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, mongo.ErrUserNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(ctx *gin.Context) {
	user := &models.User{}
	if err := ctx.BindJSON(user); err != nil {
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
		if errors.Is(err, mongo.ErrUserAlreadyExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusCreated, savedUser)
}

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
		if errors.Is(err, mongo.ErrUserNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.Is(err, mongo.ErrUserAlreadyExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

func New(s *service.Service) *UserHandler {
	return &UserHandler{service: s}
}
