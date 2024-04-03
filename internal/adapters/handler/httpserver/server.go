package httpserver

import (
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"net/http"
)

func NewServer(cfg *config.HTTP, serverHandlers []HttpHandlers) *http.Server {
	return &http.Server{
		Addr:    cfg.URL + ":" + cfg.Port,
		Handler: NewRouter(cfg.GinMode, serverHandlers),
	}
}
