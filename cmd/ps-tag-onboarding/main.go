package main

import (
	"context"
	"fmt"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/httpserver"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}
	slog.Info("Starting the application", "app", cfg.App.Name, "env", cfg.App.Env)

	db, err := repository.Connect(cfg.DB)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
	}
	if err := db.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	slog.Info("Database Connected")
	var result bson.M
	fmt.Println(db.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result))
	userRepo := repository.New(db)
	userService := service.New(userRepo)
	userHandler := httpserver.New(userService)
	router, err := httpserver.NewRouter(cfg.HTTP, *userHandler)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
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
