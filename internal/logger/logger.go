// Package logger provides a shared structured logger for the application.
// In production (APP_ENV=prod) it emits JSON; in dev it emits coloured text.
package logger

import (
	"log/slog"
	"os"
)

// L is the application-wide structured logger. Call Init() once at startup.
var L *slog.Logger

func init() {
	// Fall-back so any package that calls L before Init() never panics.
	L = slog.Default()
}

// Init configures the global logger based on the environment.
// Call this once at the very beginning of main, before any other
// service is initialised.
func Init(env string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	L = slog.New(handler)
	slog.SetDefault(L)
}
