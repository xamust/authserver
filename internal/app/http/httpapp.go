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
	"github.com/xamust/authserver/www"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"html/template"
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
	a.echoServer.GET("/health", health)

	a.echoServer.Any("/site/index.html", index)
	a.echoServer.GET("/style/*", style)
	a.echoServer.GET("/scripts/*", scripts)

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

func index(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.HTML), "html/*")
	if err != nil {
		return err
	}
	h := "http://localhost:8082"
	//if err := tmpl.Execute(c.Response().Writer, h); err != nil {
	//	return err
	//}
	return tmpl.Execute(c.Response().Writer, h)
}

func style(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.Style), "styles/*")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func scripts(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.Scripts), "scripts/*")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}
