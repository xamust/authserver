package main

import (
	"github.com/xamust/authserver/internal/app"
	"github.com/xamust/authserver/internal/config"
	"github.com/xamust/authserver/internal/xlogger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// todo: init config...
	cfg := config.MustLoad()

	// todo: init logger
	log := xlogger.NewXLogger(cfg.LogEnv)
	log.Info("starting authServer", slog.Any("cfg", cfg))
	// todo: init app

	// todo: storage....
	authApp := app.New(log, cfg.GRPCConfig, cfg.TokenTTL)

	// todo: run gRPC
	go func() {
		if err := authApp.GRPCSrv.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := authApp.HTTPSrv.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping authServer", slog.String("signal", sign.String()))
	authApp.HTTPSrv.Stop()
	authApp.GRPCSrv.Stop()
	log.Info("authServer stopped")
}
