package httpserver

import (
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"net/http"
)

func NewServer(cfg *config.HTTP, router *Router) *http.Server {
	return &http.Server{
		Addr:    cfg.URL + ":" + cfg.Port,
		Handler: router,
	}
}
