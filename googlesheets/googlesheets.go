package googlesheets

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pawlobanano/et-legacy-events-discord-bot/config"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	cfgCopy       *config.Environemnt
	sheetsService *sheets.Service
)

func Run(log *slog.Logger, cfg *config.Environemnt, client *http.Client) {
	cfgCopy = cfg
	var err error
	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Error("Unable to create Google Sheets service.", err)
	}

	http.HandleFunc("/team", AdaptReadDataHandler(getAllTeamLineups))
}

func AdaptReadDataHandler(handler func(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := func(cfgCopy *config.Environemnt) *config.Environemnt {
			return cfgCopy
		}(cfgCopy)

		handler(cfg, w, r)
	}
}

func getAllTeamLineups(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	query := cfg.GOOGLE_SHEETS_SPREADSHEET_TAB + cfg.GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE
	resp, err := sheetsService.Spreadsheets.Values.Get(cfg.GOOGLE_SHEETS_SPREADSHEET_ID, query).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resp.Values)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
