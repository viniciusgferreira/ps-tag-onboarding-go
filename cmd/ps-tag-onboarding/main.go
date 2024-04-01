package main

import (
	"context"
	"errors"
	_ "github.com/viniciusgferreira/ps-tag-onboarding-go/docs"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/handler/httpserver"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/mongodb"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/repository"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title           Tag Onboarding Go API
// @version         1.0
// @description     This is an api to manager users.

// @host      localhost:8080
func main() {
	cfg := config.New()
	slog.Info("Starting the application", "app", cfg.App.Name, "env", cfg.App.Env)
	db, err := mongodb.Connect(*cfg.DB)
	if err != nil {
		panic(err)
	}
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)

	var serverHandlers []httpserver.HttpHandlers
	serverHandlers = append(serverHandlers, httpserver.NewUserHandler(*userService))
	router := httpserver.NewRouter(cfg.HTTP.GinMode, serverHandlers)
	server := httpserver.NewServer(cfg.HTTP, router)

	var wg sync.WaitGroup
	wg.Add(1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer wg.Done()
		slog.Info("Server listening on " + server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed starting server", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-sigCh
	slog.Info("Shutting down...", "Received signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", err)
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("Server shutdown complete.")
}
