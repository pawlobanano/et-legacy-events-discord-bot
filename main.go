package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
	"github.com/pawlobanano/et-legacy-events-discord-bot/discord"
	"github.com/pawlobanano/et-legacy-events-discord-bot/googlesheets"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	config, err := config.LoadConfig(log, ".env")
	if err != nil {
		log.Error("Loading config failed.", err)
		return
	}

	discord.Run(log, config)

	googlesheets.Run(log, config, config.JwtConfig.Client(context.Background()))

	fmt.Printf("Server is listening on port %s...\n", ":8080")
	http.ListenAndServe(":8080", nil)
}
