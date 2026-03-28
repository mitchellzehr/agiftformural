package config

import "os"

type Config struct {
	Port       string
	SQLitePath string
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
	return Config{Port: port, SQLitePath: path}
}
