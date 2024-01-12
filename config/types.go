package config

import "golang.org/x/oauth2/jwt"

// Environemnt is a struct which encapsulates .env file variables.
type Environemnt struct {
	// DISCORD_BOT_API_KEY string
	JwtConfig     *jwt.Config
	SpreadsheetID string
}
