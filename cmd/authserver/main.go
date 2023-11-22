package main

import (
	"github.com/xamust/authserver/internal/app"
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

	// todo: storage....
	authApp := app.New(log, cfg.GRPCConfig.Port, cfg.GRPCConfig.Gateway, cfg.TokenTTL)

	go func() {
		if err := authApp.GRPCSrv.Run(); err != nil {
			panic(err)
		}
	}()

	if err := authApp.HTTPSrv.Run(); err != nil {
		panic(err)
	}
	// todo: run gRPC
}
