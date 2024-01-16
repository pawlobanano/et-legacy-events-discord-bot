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

const (
	TournamentName      = "DraftCup"
	TournamentEdition   = "2"
	allTeamsLineupRange = "!A:H"
	allTeamsLineupQuery = TournamentName + "#" + TournamentEdition + allTeamsLineupRange
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

	http.HandleFunc("/read", AdaptReadDataHandler(ReadData))
	http.HandleFunc("/team", AdaptReadDataHandler(getAllTeams))
}

func AdaptReadDataHandler(handler func(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := func(cfgCopy *config.Environemnt) *config.Environemnt {
			return cfgCopy
		}(cfgCopy)

		handler(cfg, w, r)
	}
}

func ReadiData(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(cfg.GOOGLE_SHEETS_SPREADSHEET_ID, allTeamsLineupQuery).Context(r.Context()).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resp.Values)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ReadData(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(cfg.GOOGLE_SHEETS_SPREADSHEET_ID, allTeamsLineupQuery).Context(r.Context()).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Data [][]interface{} `json:"data"`
	}{
		Data: resp.Values,
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getAllTeams(cfg *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(cfg.GOOGLE_SHEETS_SPREADSHEET_ID, allTeamsLineupQuery).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resp.Values)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
