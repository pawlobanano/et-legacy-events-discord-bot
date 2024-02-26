package config

import (
	"log/slog"

	"os"

	"github.com/joho/godotenv"
	"github.com/pawlobanano/tournament-discord-bot/types"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

var (
	log        *slog.Logger
	LoggingLvl = new(slog.LevelVar) // Info by default.
)

// LoadConfig loads the configuration from the .env file.
func LoadConfig(log *slog.Logger, envFilePath string) (*types.Environemnt, error) {
	err := Load(envFilePath)
	if err != nil {
		log.Error("Loading the .env file failed.", err)
		return nil, err
	}

	if err := validateEnvironmentVariables(log); err != nil {
		os.Exit(1)
	}

	envConfig := &types.Environemnt{
		DISCORD_BOT_API_KEY:                          os.Getenv("DISCORD_BOT_API_KEY"),
		ENVIRONMENT:                                  os.Getenv("ENVIRONMENT"),
		GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST:         "",
		GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL:    os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL"),
		GOOGLE_SHEETS_SPREADSHEET_ID:                 os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID"),
		GOOGLE_SHEETS_SPREADSHEET_TAB:                os.Getenv("GOOGLE_SHEETS_SPREADSHEET_TAB"),
		GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE: os.Getenv("GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE"),
		SERVER_ADDRESS:                               os.Getenv("SERVER_ADDRESS"),
		JwtConfig:                                    loadGoogleKeyJSON(),
	}

	return envConfig, nil
}

func validateEnvironmentVariables(log *slog.Logger) error {
	var err error
	if len(os.Getenv("DISCORD_BOT_API_KEY")) == 0 {
		log.Error("Environment variable DISCORD_BOT_API_KEY has not been set.")
		return err
	} else if len(os.Getenv("ENVIRONMENT")) == 0 {
		log.Error("Environment variable ENVIRONMENT has not been set.")
		return err
	} else if len(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL")) == 0 {
		log.Error("Environment variable GOOGLE_SHEETS_SPREADSHEET_ADMIN_LIST_CELL has not been set.")
		return err
	} else if len(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID")) == 0 {
		log.Error("Environment variable GOOGLE_SHEETS_SPREADSHEET_ID has not been set.")
		return err
	} else if len(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_TAB")) == 0 {
		log.Error("Environment variable GOOGLE_SHEETS_SPREADSHEET_TAB has not been set.")
		return err
	} else if len(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE")) == 0 {
		log.Error("Environment variable GOOGLE_SHEETS_SPREADSHEET_TEAM_LINEUPS_RANGE has not been set.")
		return err
	} else if len(os.Getenv("SERVER_ADDRESS")) == 0 {
		log.Error("Environment variable SERVER_ADDRESS has not been set.")
		return err
	}

	return nil
}

// loadGoogleKeyJSON utilize key.json file and returns Google's JWT Config.
func loadGoogleKeyJSON() *jwt.Config {
	creds, err := os.ReadFile("key.json")
	if err != nil {
		log.Error("Unable to read the key.json file.", err)
		os.Exit(1)
	}

	jwtConfig, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Error("Unable to create JWT config.", err)
		os.Exit(1)
	}

	return jwtConfig
}

// Load loads the environment variables from the .env file.
func Load(envFile string) error {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Error("Loading the godotenv library failed.", err)
		return err
	}

	return nil
}
