package httpserver

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

type Router struct {
	*gin.Engine
}

func NewRouter(ginMode string) *Router {
	gin.SetMode(ginMode)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &Router{router}
}

func (r *Router) SetupRouter(handlers []HttpHandlers) {
	for _, handler := range handlers {
		handler.SetupRoutes(r)
	}
}
