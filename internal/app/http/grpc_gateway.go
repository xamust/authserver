package httpapp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/xamust/authserver/pkg/authserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

func (a *App) gatewayInit() (*runtime.ServeMux, error) {
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
		return nil, fmt.Errorf("register auth grpc gateway: %w", err)
	}
	if err := authserver.RegisterConfigurationHandlerFromEndpoint(context.Background(), gw, fmt.Sprintf("%s:%d", a.conf.Host, a.conf.Port), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return nil, fmt.Errorf("register config grpc gateway: %w", err)
	}
	return gw, nil
}
