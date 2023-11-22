package app

import (
	grpcapp "github.com/xamust/authserver/internal/app/grpc"
	httpapp "github.com/xamust/authserver/internal/app/http"
	"github.com/xamust/authserver/internal/xlogger"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	HTTPSrv *httpapp.App
}

func New(log *xlogger.XLogger, grpcPort, gwGRPCPort int, tokenTTL time.Duration) *App {
	//TODO: storage init

	//TODO: init auth
	grpcApp := grpcapp.New(log, grpcPort)

	httpApp := httpapp.New(log, gwGRPCPort)

	return &App{
		GRPCSrv: grpcApp,
		HTTPSrv: httpApp,
	}
}
