package httpapp

import (
	"context"
	"fmt"
	"github.com/codemodus/swagui"
	"github.com/codemodus/swagui/suidata3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xamust/authserver/internal/xlogger"
	"github.com/xamust/authserver/pkg/authserver/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"io/fs"
	"net/http"
)

type App struct {
	log        *xlogger.XLogger
	echoServer *echo.Echo
	grpcGWPort int
}

func New(log *xlogger.XLogger, grpcGWPort int) *App {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return &App{
		log:        log,
		echoServer: e,
		grpcGWPort: grpcGWPort,
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
		// If user specified application/json in accept header,
		// we ignore google.api.Httpbody marshaler and use json
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
		}))
	//ui, err := a.swagUIInit()
	//if err != nil {
	//	return err
	//}
	a.echoServer.GET("/health", health)
	a.echoServer.Any("/", echo.WrapHandler(gw))
	//a.echoServer.Any("/swagger/*", echoSwagger.WrapHandler)
	//a.echoServer.Any("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("/api.swagger.json")))
	a.echoServer.Any("/swagger/*", echoSwagger.EchoWrapHandler(func(config *echoSwagger.Config) {
		config.URLs = []string{
			"/api.swagger.json",
		}
	}))

	a.echoServer.Any("/api.swagger.json", echo.StaticFileHandler("api.swagger.json", fs.FS(authserver.SwaggerJsonApi)))
	//a.echoServer.Any("/swaggerui/", echo.WrapHandler(http.StripPrefix("/swaggerui/", ui.Handler("/api.swagger.json"))))

	return a.echoServer.Start(fmt.Sprintf(":%d", a.grpcGWPort))
}

func (a *App) swagUIInit() (*swagui.Swagui, error) {
	return swagui.New(http.NotFoundHandler(), suidata3.New())
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
