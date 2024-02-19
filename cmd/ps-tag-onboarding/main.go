package main

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/viniciusgferreira/ps-tag-onboarding-go/docs"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/httpserver"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository/mongo"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
	"os"
)

// @title           Tag Onboarding Go API
// @version         1.0
// @description     This is an api to manager users.

// @host      localhost:8080
func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("loading environment variables", "error", err)
		os.Exit(1)
	}
	slog.Info("Starting the application", "app", cfg.App.Name, "env", cfg.App.Env)

	db := mongo.Connect(cfg.DB)

	userRepo := mongo.New(db, cfg.DB.Name)
	userService := service.New(userRepo)
	userHandler := httpserver.New(userService)
	router, err := httpserver.NewRouter(*userHandler)
	if err != nil {
		slog.Error("initializing router", "error", err)
		os.Exit(1)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf("%s:%s", cfg.HTTP.URL, cfg.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", addr)
	err = router.Serve(addr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
