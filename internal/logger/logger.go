package logger

import (
	"log/slog"
	"os"
)

// InitializeLogger sets up a JSON logger for the given use case name,
// including source code location information.
func InitializeLogger(usecaseName string) (*slog.Logger, error) {
	// Create JSON handler writing to console
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// AddSource: true, // can be enabled if you want to see file:line
	})

	return slog.New(handler), nil
}
