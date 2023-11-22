package grpcapp

import (
	"fmt"
	authgrpcserver "github.com/xamust/authserver/internal/grpc/auth"
	"github.com/xamust/authserver/internal/xlogger"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *xlogger.XLogger
	gRPCServer *grpc.Server
	port       int
}

func New(log *xlogger.XLogger, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpcserver.Register(gRPCServer)
	return &App{
		gRPCServer: gRPCServer,
		log:        log,
		port:       port,
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("starting gRPCServer on port:%d, %w", a.port, err)
	}
	a.log.With("port", a.port).With("address", l.Addr().String()).Info("gRPCServer started")
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("starting gRPCServer on port:%d, %w", a.port, err)
	}
	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
	a.log.With("port", a.port).Info("gRPCServer stopped")
}
