package auth

import (
	"context"
	authServerv1 "github.com/xamust/authserver/pkg/authserver/v1"
	"github.com/xamust/xvalidator"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authServerv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authServerv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, lr *authServerv1.LoginRequest) (*authServerv1.LoginResponse, error) {
	// todo: think about calmdown...
	if err := loginValidation(lr); err != nil {
		return nil, err
	}
	return &authServerv1.LoginResponse{Token: "token for " + lr.Email}, nil
}

func loginValidation(req *authServerv1.LoginRequest) error {
	if err := xvalidator.NewXValidator().ValidateVar(req.Email, "email"); err != nil {
		return ValidationError("email", err)
	}
	if err := xvalidator.NewXValidator().ValidateVar(req.Password, "required"); err != nil {
		return ValidationError("password", err)
	}
	return nil
}

func (s *serverAPI) Register(ctx context.Context, rr *authServerv1.RegisterRequest) (*authServerv1.RegisterResponse, error) {
	panic("not implemented")
}
