package auth

import (
	"context"
	authServerv1 "github.com/xamust/authserver/pkg/authserver/v1"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authServerv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authServerv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, lr *authServerv1.LoginRequest) (*authServerv1.LoginResponse, error) {

	return &authServerv1.LoginResponse{Token: "token for " + lr.Email}, nil
}

func (s *serverAPI) Register(ctx context.Context, rr *authServerv1.RegisterRequest) (*authServerv1.RegisterResponse, error) {
	panic("not implemented")
}
