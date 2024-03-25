package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"os"
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

func (r *Router) StartServer(url, port string, handlers []HttpHandlers) {
	for _, handler := range handlers {
		handler.SetupRoutes(r)
	}
	addr := fmt.Sprintf("%s:%s", url, port)
	err := r.Run(addr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
