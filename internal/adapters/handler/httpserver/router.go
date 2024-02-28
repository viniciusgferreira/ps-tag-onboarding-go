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

func NewRouter(userHandler UserHandler, ginMode string) *Router {
	gin.SetMode(ginMode)
	router := gin.Default()
	routes := userHandler.Routes()
	for _, r := range routes {
		router.Handle(r.Method, r.Path, r.Handler)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &Router{router}
}

func (r *Router) StartServer(url, port string) {
	addr := fmt.Sprintf("%s:%s", url, port)
	err := r.Run(addr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
