package logging

import (
	"log/slog"
	"os"
	"sync"

	"go-boilerplate/internal/config"

	"github.com/charmbracelet/log"
)

var RootLogger slog.Logger
var onceLoad sync.Once

func createLogger(cfg config.LogConfig) *slog.Logger {
	var logger *slog.Logger
	var logLevel slog.Level
	logParseError := logLevel.UnmarshalText([]byte(cfg.LogLevel))
	if logParseError != nil {
		logLevel = slog.LevelInfo
	}

	switch cfg.Style {
	case config.ColouredLogStyle:
		lvl, err := log.ParseLevel(cfg.LogLevel)
		if err != nil {
			lvl = log.InfoLevel
		}
		charmLogger := log.NewWithOptions(os.Stdout, log.Options{
			ReportTimestamp: true,
			ReportCaller:    true,
			Level:           lvl,
		})
		logger = slog.New(charmLogger)
	case config.JSONLogStyle:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		}))
	case config.PlainLogStyle:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		}))
	}

	return logger
}

func GetRootLogger(cfg config.LogConfig) *slog.Logger {
	onceLoad.Do(func() {
		RootLogger = *createLogger(cfg)
	})

	return &RootLogger
}
