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

		ch, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Error(err.Error())
		}

		if m.GuildID == "" { // It's a private message
			// The order of checking m.Content in terms of commands is important.
			if isCmdCupHelp(m) {
				s.ChannelMessageSend(ch.ID, types.DiscordMessage{Message: cupHelpCmdOutput}.Message)
				return
			}

			if isCmdCupTeams(s, m, ch) {
				s.ChannelMessageSend(
					m.ChannelID,
					types.DiscordMessage{Message: getTeamLineupsByDefaultEdition(log, cfg, s, m)}.Message)
				return
			}

			if ok, tID := isCmdCupTeamID(log, s, m, ch); ok {
				s.ChannelMessageSend(
					m.ChannelID,
					types.DiscordMessage{Message: getTeamLineupByDefaultEditionByTeamIDLetter(log, cfg, s, m, tID)}.Message)
				return
			}

			if ok, eID := isCmdCupEditionID(log, s, m, ch); ok {
				s.ChannelMessageSend(
					m.ChannelID,
					types.DiscordMessage{Message: getTeamLineupsByEditionID(log, cfg, s, m, eID)}.Message)
				return
			}

			if isCmdCupbotStatus(log, s, m) {
				s.ChannelMessageSend(ch.ID, types.DiscordMessage{Message: "Last bot's heart beat `" + DiscordSession.LastHeartbeatAck.UTC().String() + "`"}.Message)
				return
			}
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

func isCmdCupHelp(m *discordgo.MessageCreate) bool {
	if strings.EqualFold(m.Content, "!cup help") || strings.EqualFold(m.Content, "!cup h") {
		return true
	}
	return false
}

func isCmdCupTeams(s *discordgo.Session, m *discordgo.MessageCreate, ch *discordgo.Channel) bool {
	if strings.EqualFold(m.Content, "!cup teams") || strings.EqualFold(m.Content, "!cup ts") {
		return true
	}
	return false
}

func isCmdCupTeamID(log *slog.Logger, s *discordgo.Session, m *discordgo.MessageCreate, ch *discordgo.Channel) (bool, string) {
	if strings.HasPrefix(m.Content, "!cup team") || strings.HasPrefix(m.Content, "!cup t") {
		tID := extractSuffixLetter(m.Content)
		if tID == "" {
			log.Info(fmt.Sprintf("Extracting '%s' command failed.", m.Content), "value", tID)
			return false, tID
		}
		return true, tID
	}
	return false, ""
}

func isCmdCupEditionID(log *slog.Logger, s *discordgo.Session, m *discordgo.MessageCreate, ch *discordgo.Channel) (bool, int) {
	if strings.HasPrefix(m.Content, "!cup edition") || strings.HasPrefix(m.Content, "!cup e") {
		eID, err := extractSuffixNumber(m.Content)
		if err != nil || eID == 0 {
			log.Error(fmt.Sprintf("Extracting '%s' command failed.", m.Content), "error", err, "value", eID)
			return false, eID
		}
		return true, eID
	}
	return false, 0
}

func isCmdCupbotStatus(log *slog.Logger, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if strings.EqualFold(m.Content, "!cupbot status") || strings.EqualFold(m.Content, "!cupbot s") {
		log.Info(fmt.Sprintf("Cup Bot heart beat status triggered by '%s' command.", m.Content), "value", s.LastHeartbeatAck.UTC().String())
		return true
	}
	return false
}
