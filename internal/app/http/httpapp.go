package httpapp

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xamust/authserver/internal/config"
	"github.com/xamust/authserver/internal/xlogger"
	"github.com/xamust/authserver/pkg/authserver/v1"
	"github.com/xamust/xvalidator"
	"io/fs"
	"time"
)

type App struct {
	log        *xlogger.XLogger
	echoServer *echo.Echo
	conf       *config.GRPCConfig
}

func New(log *xlogger.XLogger, conf config.GRPCConfig) *App {
	e := echo.New()
	e.HideBanner = true
	e.Validator = xvalidator.NewXValidator()
	e.Use(middleware.Logger(), middleware.Recover())
	return &App{
		log:        log,
		echoServer: e,
		conf:       &conf,
	}
}

func (a *App) Run() error {
	gw, err := a.gatewayInit()
	if err != nil {
		return err
	}
	//swagger
	a.echoServer.Any("/swagger/*", echoSwagger.EchoWrapHandler(func(config *echoSwagger.Config) {
		config.URLs = []string{
			"/api.swagger.json",
		}
	}))
	a.echoServer.Any("/api.swagger.json", echo.StaticFileHandler("api.swagger.json", fs.FS(authserver.SwaggerJsonApi)))

	a.echoServer.GET("/health", health)

	// site links
	a.echoServer.Any("/site/login.html", login)
	a.echoServer.Any("/site/index.html", index)
	a.echoServer.GET("/style/*", style)
	a.echoServer.GET("/scripts/*", scripts)

	//for api
	a.echoServer.Any("/*", echo.WrapHandler(gw))

	a.log.With("port", a.conf.Gateway).With("address", a.conf.Host).Info("HTTPServer started")

	return a.echoServer.Start(fmt.Sprintf("%s:%d", a.conf.Host, a.conf.Gateway))
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := a.echoServer.Shutdown(ctx); err != nil {
		panic(err)
	}
	a.log.With("port", a.conf.Gateway).Info("HTTPserver stopped")
}
