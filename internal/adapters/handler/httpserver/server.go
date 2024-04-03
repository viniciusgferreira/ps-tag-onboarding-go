package httpserver

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"net/http"
)

type Router struct {
	*gin.Engine
}

func NewServer(cfg *config.HTTP, serverHandlers []HttpHandlers) *http.Server {
	return &http.Server{
		Addr:    cfg.URL + ":" + cfg.Port,
		Handler: NewRouter(cfg.GinMode, serverHandlers),
	}
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
