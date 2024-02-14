package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

type Router struct {
	*gin.Engine
}

func NewRouter(c *config.HTTP, userHandler UserHandler) (*Router, error) {
	router := gin.Default()
	routes := userHandler.Routes()
	for _, r := range routes {
		router.Handle(r.Method, r.Path, r.Handler)
	}
	return &Router{router}, nil
}

func (r *Router) Serve(addr string) error {
	return r.Run(addr)
}
