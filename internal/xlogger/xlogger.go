package xlogger

import (
	"log/slog"
	"os"
)

type LogEnv string

const (
	envLocal LogEnv = "local"
	envDev   LogEnv = "dev"
)

type XLogger struct {
	*slog.Logger
}

func NewXLogger(env LogEnv) *XLogger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = initNewLogger(slog.LevelDebug).With("local_env", "xamust_authServer")
	case envDev:
		log = initNewLogger(slog.LevelDebug)
	}
	return &XLogger{
		log,
	}
}

func initNewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
