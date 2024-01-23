package types

import (
	"golang.org/x/oauth2/jwt"
)

// DiscordMessage struct is used for Discord bot API messages.
type DiscordMessage struct {
	Message string
}

// Environemnt is a struct which encapsulates .env file variables.
type Environemnt struct {
	DISCORD_BOT_API_KEY                          string
	ENVIRONMENT                                  string
	GOOGLE_SHEETS_SPREADSHEET_ID                 string
	GOOGLE_SHEETS_SPREADSHEET_TAB                string
	GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE string
	SERVER_ADDRESS                               string
	JwtConfig                                    *jwt.Config
}

// GSheetsJSONResponse is a type of Google Sheets API JSON response body.
type GSheetsJSONResponse [][]string
