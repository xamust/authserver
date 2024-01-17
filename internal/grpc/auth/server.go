package auth

import (
	"context"
	"github.com/xamust/authserver/internal/validator"
	authServerv1 "github.com/xamust/authserver/pkg/authserver/v1"
	"github.com/xamust/xvalidator"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authServerv1.UnimplementedAuthServer
	authServerv1.UnimplementedConfigurationServer
}

func Register(gRPC *grpc.Server) {
	authServerv1.RegisterAuthServer(gRPC, &serverAPI{})
	authServerv1.RegisterConfigurationServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, lr *authServerv1.LoginRequest) (*authServerv1.LoginResponse, error) {
	// todo: think about calmdown...
	if err := loginValidation(lr); err != nil {
		return nil, err
	}
	return &authServerv1.LoginResponse{Token: "token for " + lr.Email}, nil
}

func (s *serverAPI) Register(ctx context.Context, rr *authServerv1.RegisterRequest) (*authServerv1.RegisterResponse, error) {
	panic("not implemented")
}

func loginValidation(req *authServerv1.LoginRequest) error {
	if err := validator.NewValidator().ValidateVar(
		[]xvalidator.InputValData{
			{Key: "email", ValData: req.Email},
			{Key: "custom_password", ValData: req.Password}}...); err != nil {
		return ValidationError("login", err)
	}
	return nil
}

func (s *serverAPI) GetConfig(ctx context.Context, in *authServerv1.GetConfigRequest) (*authServerv1.GetConfigResponse, error) {
	return &authServerv1.GetConfigResponse{Result: "result " + in.Input}, nil
}
