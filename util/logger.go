package util

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds the logger configuration read from the TOML file.
type Config struct {
	LogFile string `toml:"log_file"` // Path to the log file
	Verbose bool   `toml:"verbose"`  // Verbose mode flag (if true, logs also go to console in friendly format)
}

// InitializeLogger configures the global logger using Zerolog.
// It reads the configuration from the TOML file at ConfigFilePath.
func InitializeLogger() {
	// Open the TOML configuration file.
	configFile, err := os.Open(ConfigFilePath)
	if err != nil {
		// Fallback: use console output if the configuration file cannot be read.
		zerolog.TimeFieldFormat = DefaultTimeFormat
		log.Logger = log.Output(os.Stderr)
		return
	}
	defer configFile.Close()

	// Decode the TOML configuration.
	var cfg Config
	decoder := toml.NewDecoder(configFile)
	if _, err := decoder.Decode(&cfg); err != nil {
		zerolog.TimeFieldFormat = DefaultTimeFormat
		log.Logger = log.Output(os.Stderr)
		return
	}

	// Define um formato de tempo mais legível para o Zerolog.
	zerolog.TimeFieldFormat = time.RFC3339

	// Define o nível global de log com base no verbose.
	if cfg.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Abre ou cria o arquivo de log definido no TOML.
	logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Error().Err(err).Msg(FallbackConsoleMsg)
		log.Logger = log.Output(os.Stderr)
		return
	}

	// Se verbose, criamos um ConsoleWriter para saída amigável no console
	// e combinamos com o arquivo usando MultiLevelWriter.
	if cfg.Verbose {
		// Create a ConsoleWriter with fatih/color to format field names in cyan.
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    false,
			FormatFieldName: func(i interface{}) string {
				// Use fatih/color to print the field name in cyan.
				return color.New(color.FgCyan).Sprint(i)
			},
			FormatFieldValue: func(i interface{}) string {
				// Você pode customizar o valor se desejar; aqui usamos o valor padrão.
				return fmt.Sprintf("%s", i)
			},
		}

		// Combine a saída amigável do console com o arquivo de log.
		multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
		log.Logger = zerolog.New(multi).
			Level(zerolog.GlobalLevel()).
			With().
			Timestamp().
			Logger()
	} else {
		// Modo não-verbose: somente log para o arquivo.
		log.Logger = zerolog.New(logFile).
			Level(zerolog.GlobalLevel()).
			With().
			Timestamp().
			Logger()
	}

	log.Debug().Msg("Logger initialized")
}
