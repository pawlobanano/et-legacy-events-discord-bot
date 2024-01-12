package config

import (
	"fmt"
	"log"
	"log/slog"

	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// LoadConfig loads the configuration from the .env file.
func LoadConfig(slog *slog.Logger, envFilePath string) (*Environemnt, error) {
	err := Load(envFilePath)
	if err != nil {
		slog.Error("Loading .env file.", err)
		return nil, err
	}

	creds, err := os.ReadFile("./key.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	jwtConfig, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to create JWT config: %v", err)
	}

	// DISCORD_BOT_API_KEY := os.Getenv("DISCORD_BOT_API_KEY")

	// if len(DISCORD_BOT_API_KEY) != 32 {
	// 	slog.Error(fmt.Sprintf("%s %d", "DISCORD_BOT_API_KEY must be equal to 32 characters. But was", len(DISCORD_BOT_API_KEY)))
	// }

	config := &Environemnt{
		// DISCORD_BOT_API_KEY: DISCORD_BOT_API_KEY,
		JwtConfig:     jwtConfig,
		SpreadsheetID: os.Getenv("GOOGLE_SPREADSHEET_ID"),
	}

	return config, nil
}

// Load loads the environment variables from the .env file.
func Load(envFile string) error { // Solution to differentiate .env file path for unit or benchmark tests; source: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
	err := godotenv.Load(dir(envFile))
	if err != nil {
		slog.Error(err.Error())
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
