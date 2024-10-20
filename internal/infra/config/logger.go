package config

import (
	"log/slog"
	"os"
)

func ConfigureLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
