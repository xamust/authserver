package app

import (
	grpcapp "github.com/xamust/authserver/internal/app/grpc"
	httpapp "github.com/xamust/authserver/internal/app/http"
	"github.com/xamust/authserver/internal/config"
	"github.com/xamust/authserver/internal/xlogger"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	HTTPSrv *httpapp.App
}

func New(log *xlogger.XLogger, conf config.GRPCConfig, tokenTTL time.Duration) *App {
	//TODO: storage init

	//TODO: init auth
	grpcApp := grpcapp.New(log, conf.Port)

	httpApp := httpapp.New(log, conf)

	return &App{
		GRPCSrv: grpcApp,
		HTTPSrv: httpApp,
	}
}
