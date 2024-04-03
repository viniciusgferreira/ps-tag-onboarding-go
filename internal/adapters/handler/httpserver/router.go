package httpserver

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(ginMode string, handlers []HttpHandlers) *Router {
	gin.SetMode(ginMode)
	router := &Router{gin.Default()}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, handler := range handlers {
		handler.SetupRoutes(router)
	}
	return router
}

func (r *Router) SetupRouter(handlers []HttpHandlers) {
	for _, handler := range handlers {
		handler.SetupRoutes(r)
	}
}
