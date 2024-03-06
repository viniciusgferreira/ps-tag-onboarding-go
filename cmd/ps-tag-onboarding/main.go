package main

import (
	_ "github.com/viniciusgferreira/ps-tag-onboarding-go/docs"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/httpserver"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/mongodb"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
)

// @title           Tag Onboarding Go API
// @version         1.0
// @description     This is an api to manager users.

// @host      localhost:8080
func main() {
	cfg := config.New()
	slog.Info("Starting the application", "app", cfg.App.Name, "env", cfg.App.Env)
	db := mongodb.Connect(*cfg.DB)
	userRepo := repository.New(db)
	userService := service.New(userRepo)
	userHandler := httpserver.New(userService)
	router := httpserver.NewRouter(*userHandler, cfg.HTTP.GinMode)

	router.StartServer(cfg.HTTP.URL, cfg.HTTP.Port)
}
