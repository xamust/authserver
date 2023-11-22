package httpapp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xamust/authserver/internal/config"
	"github.com/xamust/authserver/internal/xlogger"
	"github.com/xamust/authserver/pkg/authserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"io/fs"
	"net/http"
	"time"
)

type App struct {
	log        *xlogger.XLogger
	echoServer *echo.Echo
	conf       *config.GRPCConfig
}

func New(log *xlogger.XLogger, conf config.GRPCConfig) *App {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return &App{
		log:        log,
		echoServer: e,
		conf:       &conf,
	}
}

func (a *App) Run() error {
	gw := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
			},
		}),
		//If user specified application/json in accept header,
		//we ignore google.api.Httpbody marshaler and use json
		runtime.WithMarshalerOption("application/json", &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "X-Request-Id", "X-Real-Ip", "X-Forwarded-Host", "X-Visitor-Id":
				return key, true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
		runtime.WithMetadata(func(c context.Context, r *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				"origin-url": r.URL.String(),
			})
		}),
	)

	if err := authserver.RegisterAuthHandlerFromEndpoint(context.Background(), gw, fmt.Sprintf("%s:%d", a.conf.Host, a.conf.Port), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return fmt.Errorf("register grpc gateway: %w", err)
	}
	a.echoServer.Any("/swagger/*", echoSwagger.EchoWrapHandler(func(config *echoSwagger.Config) {
		config.URLs = []string{
			"/api.swagger.json",
		}
	}))
	a.echoServer.GET("/", health)
	a.echoServer.Any("/*", echo.WrapHandler(gw))
	a.echoServer.Any("/api.swagger.json", echo.StaticFileHandler("api.swagger.json", fs.FS(authserver.SwaggerJsonApi)))

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

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
