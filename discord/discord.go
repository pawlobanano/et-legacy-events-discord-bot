package discord

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pawlobanano/et-legacy-events-discord-bot/types"
)

var DiscordSession *discordgo.Session

const cupHelpCmdOutput = "`!cup edition <number> | e <num>`\n`!cup help | h`\n`!cup team <letter> | t <let>`\n`!cup teams | ts`\n"

func Run(log *slog.Logger, cfg *types.Environemnt) {
	DiscordSession, err := discordgo.New("Bot " + cfg.DISCORD_BOT_API_KEY)
	if err != nil {
		log.Error("Error instantiating bot.")
		os.Exit(1)
	}

	DiscordSession.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Error(err.Error())
		}

		// The order of checking m.Content in terms of commands is important.
		if strings.EqualFold(m.Content, "!cup help") || strings.EqualFold(m.Content, "!cup h") {
			s.ChannelMessageSend(channel.ID, types.DiscordMessage{Message: cupHelpCmdOutput}.Message)
			return
		}

		if strings.EqualFold(m.Content, "!cup teams") || strings.EqualFold(m.Content, "!cup ts") {
			s.ChannelMessageSend(
				m.ChannelID,
				types.DiscordMessage{Message: getTeamLineupsByDefaultEdition(log, cfg, s, m)}.Message)
			return
		}

		if strings.HasPrefix(m.Content, "!cup team") || strings.HasPrefix(m.Content, "!cup t") {
			tID := extractSuffixLetter(m.Content)
			if tID == "" {
				log.Info(fmt.Sprintf("Extracting '%s' command failed.", m.Content), "value", tID)
				return
			}
			s.ChannelMessageSend(
				m.ChannelID,
				types.DiscordMessage{Message: getTeamLineupByDefaultEditionByTeamIDLetter(log, cfg, s, m, tID)}.Message)
			return
		}

		if strings.HasPrefix(m.Content, "!cup edition") || strings.HasPrefix(m.Content, "!cup e") {
			eID, err := extractSuffixNumber(m.Content)
			if err != nil || eID == 0 {
				log.Error(fmt.Sprintf("Extracting '%s' command failed.", m.Content), "error", err, "value", eID)
				return
			}
			s.ChannelMessageSend(
				m.ChannelID,
				types.DiscordMessage{Message: getTeamLineupsByEditionID(log, cfg, s, m, eID)}.Message)
			return
		}

		if strings.EqualFold(m.Content, "!cupbot status") || strings.EqualFold(m.Content, "!cupbot s") {
			s.ChannelMessageSend(channel.ID, types.DiscordMessage{Message: "Last bot's heart beat `" + DiscordSession.LastHeartbeatAck.UTC().String() + "`"}.Message)
			return
		}
	})

	err = DiscordSession.Open()
	if err != nil {
		log.Error("Failed to run session.Open() func which creates a websocket connection to Discord.", err)
		os.Exit(1)
	}
	defer DiscordSession.Close()

	log.Info("Bot is now running.")
	log.Info("Press Ctrl+C to exit.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
