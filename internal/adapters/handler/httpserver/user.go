package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/ports"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"net/http"
)

type UserHandler struct {
	service ports.UserService
}

func (h *UserHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/users", Handler: h.GetByID},
		{Method: http.MethodPost, Path: "/users", Handler: h.Create},
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, "user")
}

func (h *UserHandler) Create(c *gin.Context) {
	c.IndentedJSON(201, "test:test")
}

func New(s *service.Service) *UserHandler {
	return &UserHandler{service: s}
}
