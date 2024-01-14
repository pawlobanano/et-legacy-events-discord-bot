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
	readRange = "DraftCup#2!A:H"
)

var (
	envConfigCopy *config.Environemnt
	sheetsService *sheets.Service
)

func Run(log *slog.Logger, envConfig *config.Environemnt, client *http.Client) {
	envConfigCopy = envConfig
	var err error
	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Error("Unable to create Google Sheets service: ", err)
	}

	http.HandleFunc("/read", AdaptReadDataHandler(ReadData))
}

func AdaptReadDataHandler(handler func(envConfig *config.Environemnt, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envConfig := func(envConfigCopy *config.Environemnt) *config.Environemnt {
			return envConfigCopy
		}(envConfigCopy)

		handler(envConfig, w, r)
	}
}

func ReadiData(envConfig *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(envConfig.GOOGLE_SHEETS_SPREADSHEET_ID, readRange).Context(r.Context()).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resp.Values)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ReadData(envConfig *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(envConfig.GOOGLE_SHEETS_SPREADSHEET_ID, readRange).Context(r.Context()).Do()
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
