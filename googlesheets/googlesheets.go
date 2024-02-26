package googlesheets

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pawlobanano/tournament-discord-bot/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var SheetsService *sheets.Service

func Run(log *slog.Logger, cfg *types.Environemnt) (err error) {
	SheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(cfg.JwtConfig.Client(context.Background())))
	if err != nil {
		log.Error("Unable to create Google Sheets service.", err)
		return err
	}

	cfg.GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST = GetAdmins(log, cfg)

	return nil
}

func GetAdmins(log *slog.Logger, cfg *types.Environemnt) string {
	resp, err :=
		SheetsService.Spreadsheets.Values.Get(
			cfg.GOOGLE_SHEETS_SPREADSHEET_ID,
			cfg.GOOGLE_SHEETS_SPREADSHEET_TAB+cfg.GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL,
		).Do()
	if err != nil {
		log.Error(fmt.Sprintf("Retriving admin list cell '%s' from Google Sheet Tab '%s' failed.", cfg.GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL, cfg.GOOGLE_SHEETS_SPREADSHEET_TAB))
		return ""
	}

	if len(resp.Values) == 0 {
		log.Info(fmt.Sprintf("Retriving admin list cell '%s' from Google Sheet Tab '%s' returned no data.", cfg.GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL, cfg.GOOGLE_SHEETS_SPREADSHEET_TAB))
		return ""
	}

	return fmt.Sprintf("%s", resp.Values[0][0])
}
