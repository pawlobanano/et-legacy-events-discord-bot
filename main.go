package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pawlobanano/et-legacy-events-discord-bot/api"
	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
	"github.com/pawlobanano/et-legacy-events-discord-bot/discord"
	"github.com/pawlobanano/et-legacy-events-discord-bot/googlesheets"
)

var (
	slogHandlerOptions = &slog.HandlerOptions{AddSource: false, Level: config.LoggingLvl}
	wg                 sync.WaitGroup
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log := slog.New(slog.NewJSONHandler(os.Stdout, slogHandlerOptions))

	cfg, err := config.LoadConfig(log, ".env")
	if err != nil {
		log.LogAttrs(context.Background(), slog.LevelError, "Loading the config package failed.", slog.String("msg", err.Error())) // As an optimized example.
		return
	}

	if cfg.ENVIRONMENT == "local" {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	}

	if err = googlesheets.Run(log, cfg); err != nil {
		log.Error("Unable to create Google Sheets service.", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		discord.Run(log, cfg)
	}()

	server := api.NewServer(cfg, log)

	go func() {
		log.Info("HTTP server is running.", "server address", cfg.SERVER_ADDRESS)
		if err := server.Start(); err != nil {
			log.Error("HTTP server error.", "msg", err)
			os.Exit(1)
		}
	}()

	<-interrupt
	log.Info("Received interrupt signal. Shutting down...")

	wg.Wait()
	log.Info("Shutdown complete.")

	os.Exit(0)
}
