package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	SQLitePath  string
	MuralURL    string
	MuralAPIKey string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	path := os.Getenv("SQLITE_PATH")
	if path == "" {
		path = "app.db"
	}
	muralURL := strings.TrimSpace(os.Getenv("MURAL_BASE_URL"))
	if muralURL == "" {
		muralURL = "https://api-staging.muralpay.com"
	}
	apiKey := strings.TrimSpace(os.Getenv("MURAL_API_KEY"))
	return Config{
		Port:        port,
		SQLitePath:  path,
		MuralURL:    muralURL,
		MuralAPIKey: apiKey,
	}
}
