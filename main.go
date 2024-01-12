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

var (
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wg  sync.WaitGroup
)

func main() {
	config, err := config.LoadConfig(log, ".env")
	if err != nil {
		log.Error("Loading config failed.", err)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		discord.Run(log, config)
	}()

	googlesheets.Run(log, config, config.JwtConfig.Client(context.Background()))

	go func() {
		log.Info("Server is listening on port :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Error("HTTP server error: ", err)
		}
	}()

	<-interrupt
	log.Info("Received interrupt signal. Shutting down...")

	wg.Wait()
	log.Info("Shutdown complete.")
}
