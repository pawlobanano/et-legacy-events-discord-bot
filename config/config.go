package config

import (
	"fmt"
	"log/slog"

	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

var (
	log        *slog.Logger
	LoggingLvl = new(slog.LevelVar) // Info by default.
)

// LoadConfig loads the configuration from the .env file.
func LoadConfig(log *slog.Logger, envFilePath string) (*Environemnt, error) {
	err := Load(envFilePath)
	if err != nil {
		log.Error("Loading .env file.", err)
		return nil, err
	}

	if len(os.Getenv("DISCORD_BOT_API_KEY")) == 0 {
		log.Error("Environment variable DISCORD_BOT_API_KEY has not been set.")
		os.Exit(1)
	}

	if len(os.Getenv("ENVIRONMENT")) == 0 {
		log.Error("Environment variable ENVIRONMENT has not been set.")
		os.Exit(1)
	}

	if len(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID")) == 0 {
		log.Error("Environment variable GOOGLE_SHEETS_SPREADSHEET_ID has not been set.")
		os.Exit(1)
	}

	config := &Environemnt{
		DISCORD_BOT_API_KEY:          os.Getenv("DISCORD_BOT_API_KEY"),
		ENVIRONMENT:                  os.Getenv("ENVIRONMENT"),
		GOOGLE_SHEETS_SPREADSHEET_ID: os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID"),
		JwtConfig:                    loadGoogleKeyJSON(),
	}

	return config, nil
}

func loadGoogleKeyJSON() *jwt.Config {
	creds, err := os.ReadFile("./key.json")
	if err != nil {
		log.Error("Unable to read key.json file.", err)
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
func Load(envFile string) error { // Solution to differentiate .env file path for unit or benchmark tests; source: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
	err := godotenv.Load(dir(envFile))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
