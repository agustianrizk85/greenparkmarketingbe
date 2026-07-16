// Command server starts the Marketing "Qualified Demand Control Tower" API.
//
// It wires the layers together — repository -> service -> HTTP transport — and
// runs an HTTP server with graceful shutdown.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"greenpark/marketing/internal/authmw"
	"greenpark/marketing/internal/config"
	"greenpark/marketing/internal/repository"
	"greenpark/marketing/internal/service"
	httptransport "greenpark/marketing/internal/transport/http"
)

func main() {
	cfg := config.Load()

	// Dependency wiring (composition root). Use PostgreSQL when a DSN is set,
	// otherwise persist to the JSON file.
	var (
		repo repository.MarketingRepository
		err  error
	)
	if cfg.DatabaseURL != "" {
		repo, err = repository.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("marketing: postgres: %v", err)
		}
		log.Println("marketing: using PostgreSQL store")
	} else {
		repo, err = repository.NewRepository(cfg.DataPath)
		if err != nil {
			log.Fatalf("marketing: failed to open data store %q: %v", cfg.DataPath, err)
		}
		log.Println("marketing: using file store")
	}
	svc := service.New(repo)
	verifier, err := authmw.New(authmw.Options{
		JWKSURL:    cfg.AuthJWKSURL,
		Department: "marketing",
		Issuer:     cfg.AuthIssuer,
	})
	if err != nil {
		log.Fatalf("marketing: auth verifier: %v", err)
	}
	handler := httptransport.NewHandler(svc, verifier)
	router := httptransport.NewRouter(handler, cfg.AllowOrigin)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in a goroutine so main can wait for shutdown signals.
	go func() {
		log.Printf("marketing API listening on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("marketing: server error: %v", err)
		}
	}()

	// Wait for an interrupt signal for graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("marketing: shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("marketing: graceful shutdown failed: %v", err)
	}
	log.Println("marketing: stopped")
}
