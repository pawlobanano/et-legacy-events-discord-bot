package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
	"github.com/pawlobanano/et-legacy-events-discord-bot/discord"
	"github.com/pawlobanano/et-legacy-events-discord-bot/googlesheets"
)

var wg sync.WaitGroup

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LoggingLvl}))

	cfg, err := config.LoadConfig(log, ".env")
	if err != nil {
		log.LogAttrs(context.Background(), slog.LevelError, "Loading the config package failed.", slog.String("msg", err.Error())) // As an optimize example.
		return
	}

	if cfg.ENVIRONMENT == "local" {
		config.LoggingLvl.Set(slog.LevelDebug)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		discord.Run(log, cfg)
	}()

	googlesheets.Run(log, cfg, cfg.JwtConfig.Client(context.Background()))

	go func() {
		log.Info("Server is listening on port :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Error("HTTP server error.", "msg", err)
		}
	}()

	<-interrupt
	log.Info("Received interrupt signal. Shutting down...")
	wg.Wait()
	log.Info("Shutdown complete.")
}
