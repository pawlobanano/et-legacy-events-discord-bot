package config

import "golang.org/x/oauth2/jwt"

// Environemnt is a struct which encapsulates .env file variables.
type Environemnt struct {
	DISCORD_BOT_API_KEY          string
	ENVIRONMENT                  string
	GOOGLE_SHEETS_SPREADSHEET_ID string
	JwtConfig                    *jwt.Config
}

// Response is a type of HTTP response body.
type Response [][]string
