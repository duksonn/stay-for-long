package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/duksonn/stay-for-long/cmd/config"
	"github.com/duksonn/stay-for-long/cmd/di"
	internalhttp "github.com/duksonn/stay-for-long/internal/infra/http"
)

func main() {
	cfg := config.Load()

	deps := di.Init()
	log.Printf("Dependencies init successfully")

	router, err := internalhttp.Routes(deps)
	if err != nil {
		log.Fatalf("Could not start router: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Printf("Starting server on %v", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not start server: %v", err)
	}
}
