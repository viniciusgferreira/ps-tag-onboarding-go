package main

import (
	"fmt"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/httpserver"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository/mongo"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("loading environment variables", "error", err)
		os.Exit(1)
	}
	slog.Info("Starting the application", "app", cfg.App.Name, "env", cfg.App.Env)

	db := mongo.Connect(cfg.DB)

	userRepo := mongo.New(db)
	userService := service.New(userRepo)
	userHandler := httpserver.New(userService)
	router, err := httpserver.NewRouter(*userHandler)
	if err != nil {
		slog.Error("initializing router", "error", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%s", cfg.HTTP.URL, cfg.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", addr)
	err = router.Serve(addr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
