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
	configCopy    *config.Environemnt
	sheetsService *sheets.Service
)

func Run(log *slog.Logger, config *config.Environemnt, client *http.Client) {
	configCopy = config
	var err error
	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Error("Unable to create Google Sheets service: ", err)
	}

	http.HandleFunc("/read", AdaptReadDataHandler(ReadData))
}

func AdaptReadDataHandler(handler func(config *config.Environemnt, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config := func(configCopy *config.Environemnt) *config.Environemnt {
			return configCopy
		}(configCopy)

		handler(config, w, r)
	}
}

func ReadiData(config *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(config.SpreadsheetID, readRange).Context(r.Context()).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resp.Values)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ReadData(config *config.Environemnt, w http.ResponseWriter, r *http.Request) {
	resp, err := sheetsService.Spreadsheets.Values.Get(config.SpreadsheetID, readRange).Context(r.Context()).Do()
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
