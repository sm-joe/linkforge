package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sm-joe/linkforge/internal/config"
	httpserver "github.com/sm-joe/linkforge/internal/http"
	"github.com/sm-joe/linkforge/internal/link"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	cfg := config.Load()
	logger := httpserver.NewLogger()

	db, err := link.OpenDatabase(cfg.DatabasePath)
	if err != nil {
		logger.Error(
			"database initialization failed",
			"error",
			err,
		)
		os.Exit(1)
	}
	defer db.Close()

	repository := link.NewSQLiteRepository(db)
	linkService := link.NewService(repository)

	linkHandler := httpserver.NewLinkHandler(linkService)
	listHandler := httpserver.NewListHandler(linkService)
	managementHandler := httpserver.NewLinkManagementHandler(linkService)
	redirectHandler := httpserver.NewRedirectHandler(linkService)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/healthz",
		func(w http.ResponseWriter, r *http.Request) {
			writeHealthJSON(w)
		},
	)

	// Create link + list links
	mux.HandleFunc(
		"/api/v1/links",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {

			case http.MethodGet:
				listHandler.List(w, r)

			case http.MethodPost:
				linkHandler.Create(w, r)

			default:
				http.NotFound(w, r)
			}
		},
	)

	// Manage individual links
	mux.HandleFunc(
		"/api/v1/links/",
		func(w http.ResponseWriter, r *http.Request) {
			switch {

			case r.Method == http.MethodGet:
				managementHandler.Get(w, r)

			case r.Method == http.MethodDelete:
				managementHandler.Delete(w, r)

			case r.Method == http.MethodPost &&
				strings.HasSuffix(r.URL.Path, "/disable"):
				managementHandler.Disable(w, r)

			case r.Method == http.MethodPost &&
				strings.HasSuffix(r.URL.Path, "/enable"):
				managementHandler.Enable(w, r)

			default:
				http.NotFound(w, r)
			}
		},
	)

	// Public short URL redirect
	mux.HandleFunc(
		"/",
		redirectHandler.Redirect,
	)

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: httpserver.WithCORS(
			mux,
		),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(
			"LinkForge API starting",
			"address",
			server.Addr,
			"database",
			cfg.DatabasePath,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {

	case err := <-serverErrors:
		logger.Error(
			"server failed",
			"error",
			err,
		)
		os.Exit(1)

	case sig := <-shutdown:
		logger.Info(
			"shutdown signal received",
			"signal",
			sig.String(),
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(
			"graceful shutdown failed",
			"error",
			err,
		)
		os.Exit(1)
	}

	logger.Info("LinkForge API stopped")
}

func writeHealthJSON(
	w http.ResponseWriter,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(
		healthResponse{
			Status: "ok",
		},
	)
}
