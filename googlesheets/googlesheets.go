package googlesheets

import (
	"context"
	"log/slog"

	"github.com/pawlobanano/et-legacy-events-discord-bot/types"
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

	return nil
}
