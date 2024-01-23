package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pawlobanano/et-legacy-events-discord-bot/types"
)

func getTeamLineupsByDefaultEdition(log *slog.Logger, cfg *types.Environemnt, s *discordgo.Session, m *discordgo.MessageCreate) (outputMessage string) {
	resp, err := http.DefaultClient.Get("http://" + cfg.SERVER_ADDRESS + "/cup")
	if err != nil {
		log.Error(fmt.Sprintf("GET http://%s/cup failed", cfg.SERVER_ADDRESS), err.Error(), slog.Int("HTTP status code", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read the response body.", err)
	}

	var data types.GSheetsJSONResponse
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Error("Can not unmarshal JSON.", err)
	}

	return createDiscordBotMultilineMessage(cfg, data, cfg.GOOGLE_SHEETS_SPREADSHEET_TAB, log).format()
}

func getTeamLineupsByEditionID(log *slog.Logger, cfg *types.Environemnt, s *discordgo.Session, m *discordgo.MessageCreate, eID int) (outputMessage string) {
	resp, err := http.DefaultClient.Get("http://" + cfg.SERVER_ADDRESS + "/cup/edition?id=" + fmt.Sprint(eID))
	if err != nil {
		log.Error(fmt.Sprintf("GET http://%s/cup/edition?id=%d failed", cfg.SERVER_ADDRESS, eID), err.Error(), slog.Int("HTTP status code", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read the response body.", err)
	}

	var data types.GSheetsJSONResponse
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Error("Can not unmarshal JSON.", err)
	}

	edition := regexp.MustCompile(`\d+`).ReplaceAllString(cfg.GOOGLE_SHEETS_SPREADSHEET_TAB, "") + fmt.Sprint(eID)

	return createDiscordBotMultilineMessage(cfg, data, edition, log).format()
}

func getTeamLineupByDefaultEditionByTeamIDLetter(log *slog.Logger, cfg *types.Environemnt, s *discordgo.Session, m *discordgo.MessageCreate, tID string) (outputMessage string) {
	resp, err := http.DefaultClient.Get("http://" + cfg.SERVER_ADDRESS + "/cup/team?id=" + tID)
	if err != nil {
		log.Error(fmt.Sprintf("GET http://%s/cup/team?id=%s failed", cfg.SERVER_ADDRESS, tID), err.Error(), slog.Int("HTTP status code", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read the response body.", err)
	}

	if string(jsonData) == "" { // Failed to find a team by ID letter. Logged by API handler.
		return
	}

	var data types.GSheetsJSONResponse
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Error("Can not unmarshal JSON.", err)
	}

	return createDiscordBotMultilineMessage(cfg, data, cfg.GOOGLE_SHEETS_SPREADSHEET_TAB, log).format()
}

func createDiscordBotMultilineMessage(cfg *types.Environemnt, data types.GSheetsJSONResponse, edition string, log *slog.Logger) *multilineString {
	multilineStr := newMultilineString("# " + edition + " | Team lineups")
	for _, row := range data {
		for j, value := range row {
			parts := strings.Split(value, " ")
			if len(parts) == 2 {
				countryCode := parts[1]
				if j == 0 {
					multilineStr.append("## " + value) // Team
					continue
				} else if j == 1 && value != "" {
					multilineStr.append("- :flag_" + countryCode + ": **" + parts[0] + "**") // Captain
					continue
				}
				multilineStr.append("- :flag_" + countryCode + ": " + parts[0]) // Standard player
			} else {
				log.Info("Wrong player name input string.", "Value", value)
				multilineStr.append("- " + parts[0]) // Standard player
			}
		}
		multilineStr.append("")
	}
	multilineStr.append("`!cup help` _for available commands_")

	return multilineStr
}
