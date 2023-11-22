package main

import (
	"github.com/xamust/authserver/internal/config"
	"github.com/xamust/authserver/internal/xlogger"
	"log/slog"
)

func main() {
	// todo: init config...
	cfg := config.MustLoad()

	// todo: init logger
	log := xlogger.NewXLogger(cfg.LogEnv)
	log.Info("starting authServer", slog.Any("cfg", cfg))
	// todo: init app

	// todo: run gRPC
}
