package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pawlobanano/tournament-discord-bot/googlesheets"
)

func (s *Server) getTeamLineupsByDefaultEdition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err :=
		googlesheets.SheetsService.Spreadsheets.Values.Get(
			s.config.GOOGLE_SHEETS_SPREADSHEET_ID,
			s.config.GOOGLE_SHEETS_SPREADSHEET_TAB+s.config.GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE,
		).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) getTeamLineupsByEditionID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	edition := regexp.MustCompile(`\d+`).ReplaceAllString(s.config.GOOGLE_SHEETS_SPREADSHEET_TAB, "") + fmt.Sprint(id)

	resp, err :=
		googlesheets.SheetsService.Spreadsheets.Values.Get(
			s.config.GOOGLE_SHEETS_SPREADSHEET_ID,
			edition+s.config.GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE,
		).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) getPlaythroughByDefaultEdition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err :=
		googlesheets.SheetsService.Spreadsheets.Values.Get(
			s.config.GOOGLE_SHEETS_SPREADSHEET_ID,
			s.config.GOOGLE_SHEETS_SPREADSHEET_TAB+"!A13:K32",
		).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) getTeamLineupByDefaultEditionByTeamIDLetter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	resp, err :=
		googlesheets.SheetsService.Spreadsheets.Values.Get(
			s.config.GOOGLE_SHEETS_SPREADSHEET_ID,
			s.config.GOOGLE_SHEETS_SPREADSHEET_TAB+fmt.Sprintf("!%s1:%s7", id, id),
		).MajorDimension("COLUMNS").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if string(jsonData) == "null" {
		s.logger.Info("Failed to find a team by ID letter.", "value", id)
		return
	}
	w.Write(jsonData)
}
