package discord

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
)

func Run(log *slog.Logger, envConfig *config.Environemnt) {
	bot, err := discordgo.New("Bot " + envConfig.DISCORD_BOT_API_KEY)
	if err != nil {
		log.Error("Error instantiating bot.")
		os.Exit(1)
	}

	bot.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		if mess.Author.ID == sess.State.User.ID {
			return
		}

		if strings.HasPrefix(mess.Content, "!cupbot status") || strings.HasPrefix(mess.Content, "!cupbot s") {
			sess.ChannelMessageSend(mess.ChannelID, "Last bot's heart beat: `"+bot.LastHeartbeatAck.UTC().String()+"`")
			return
		}
	})

	err = bot.Open()
	if err != nil {
		log.Error("Failed to run bot.Open() func which creates a websocket connection to Discord.", err)
		os.Exit(1)
	}
	defer bot.Close()

	log.Info("Bot is now running.")
	log.Info("Press Ctrl+C to exit.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
