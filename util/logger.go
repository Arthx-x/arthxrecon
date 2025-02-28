package logger

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogFile string `json:"log_file"`
	Verbose bool   `json:"verbose"`
}

// InitializeLogger configures the global logger using Zerolog.
// It reads the log file path and verbose flag from config/config.json.
func InitializeLogger() {
	file, err := os.Open("config/config.json")
	if err != nil {
		// Fallback para o console
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(os.Stderr)
		return
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(os.Stderr)
		return
	}

	// Define o nível global de log
	if cfg.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Abre o arquivo de log conforme configurado
	logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open log file, using console output")
		log.Logger = log.Output(os.Stderr)
		return
	}

	// Se verbose, também imprime no console
	if cfg.Verbose {
		log.Logger = zerolog.New(zerolog.MultiLevelWriter(logFile, os.Stderr)).
			Level(zerolog.GlobalLevel()).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(logFile).
			Level(zerolog.GlobalLevel()).With().Timestamp().Logger()
	}
	log.Debug().Msg("Logger initialized")
}
