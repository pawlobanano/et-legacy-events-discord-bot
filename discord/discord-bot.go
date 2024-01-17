package discord

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
)

func Run(log *slog.Logger, cfg *config.Environemnt) {
	session, err := discordgo.New("Bot " + cfg.DISCORD_BOT_API_KEY)
	if err != nil {
		log.Error("Error instantiating bot.")
		os.Exit(1)
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Error(err.Error())
		}

		if strings.EqualFold(m.Content, "!cup help") || strings.EqualFold(m.Content, "!cup h") {
			s.ChannelMessageSend(channel.ID, "`!cup help | h`\n`!cup teams | t`\n`!cupbot status | s`")
			return
		}

		if strings.EqualFold(m.Content, "!cup teams") || strings.EqualFold(m.Content, "!cup t") {
			getAllTeamLineups(log, cfg, s, m)
			return
		}

		if strings.EqualFold(m.Content, "!cupbot status") || strings.EqualFold(m.Content, "!cupbot s") {
			s.ChannelMessageSend(channel.ID, "Last bot's heart beat `"+session.LastHeartbeatAck.UTC().String()+"`")
			return
		}
	})

	err = session.Open()
	if err != nil {
		log.Error("Failed to run session.Open() func which creates a websocket connection to Discord.", err)
		os.Exit(1)
	}
	defer session.Close()

	log.Info("Bot is now running.")
	log.Info("Press Ctrl+C to exit.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func getAllTeamLineups(log *slog.Logger, cfg *config.Environemnt, s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.DefaultClient.Get("http://localhost:8080/team")
	if err != nil {
		log.Error("GET /team failed.", err.Error(), slog.Int("HTTP status code", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Failed to read the response body.", err)
		}

		var result config.Response
		if err := json.Unmarshal(body, &result); err != nil {
			log.Error("Can not unmarshal JSON.", err)
		}

		multilineStr := NewMultilineString("# " + cfg.GOOGLE_SHEETS_SPREADSHEET_TAB + " | Team lineups")
		for _, row := range result {
			for j, value := range row {
				parts := strings.Split(value, " ")
				if len(parts) == 2 {
					countryCode := parts[1]
					if j == 0 {
						multilineStr.Append("## " + value) // Team
						continue
					} else if j == 1 && value != "" {
						multilineStr.Append("- :flag_" + countryCode + ": **" + parts[0] + "**") // Captain
						continue
					}
					multilineStr.Append("- :flag_" + countryCode + ": " + parts[0]) // Standard player
				} else {
					log.Info("Wrong player name input string.", "Value", value)
					multilineStr.Append("- " + parts[0]) // Standard player
				}
			}
			multilineStr.Append("")
		}
		multilineStr.Append("`!cup help` _to check available commands_")

		s.ChannelMessageSend(channel.ID, multilineStr.Format())
	} else {
		log.Error("Request failed.", slog.Int("HTTP response status code", resp.StatusCode), "request URI", resp.Request.URL.String())
	}
}
