package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetLoggerConfig() {
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			AddSource: true,	
		}),
	))
}