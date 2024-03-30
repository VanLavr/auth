package logging

import (
	"log/slog"
	"os"
)

type Logger struct {
	l *slog.Logger
}

func New() *Logger {
	log := new(Logger)
	log.l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return log
}

func (l *Logger) SetAsDefault() {
	slog.SetDefault(l.l)
}
