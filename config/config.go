package config

import "os"

// Config represents the configuration of the application.
type Config struct {
	TelegramToken              string
	BinanceAPIKey              string
	BinanceAPISecret           string
	PostgreSQLConnectionString string
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() Config {
	return Config{
		TelegramToken:              os.Getenv("TELEGRAM_TOKEN"),
		BinanceAPIKey:              os.Getenv("BINANCE_API_KEY"),
		BinanceAPISecret:           os.Getenv("BINANCE_API_SECRET"),
		PostgreSQLConnectionString: os.Getenv("POSTGRESQL_CONNECTION_STRING"),
	}
}
