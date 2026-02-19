package httpserver

import (
	"net/http"
	"time"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/config"
)

func New(handler http.Handler, cfg config.HTTPConfig) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}
}
