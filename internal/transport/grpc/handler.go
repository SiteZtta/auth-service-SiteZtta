package grpc

import (
	"context"

	sitezttav1 "github.com/SiteZtta/protos-SiteZtta/gen/go/auth"
	"google.golang.org/grpc"
)

type handler struct {
	sitezttav1.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	sitezttav1.RegisterAuthServiceServer(gRPC, &handler{})
}

func (h *handler) CreateUser(ctx context.Context, req *sitezttav1.SignUpRequest) (*sitezttav1.UserIdResponse, error) {
	//panic("implement me")
	return nil, nil
}

func (h *handler) GenerateToken(ctx context.Context, req *sitezttav1.SignInRequest) (*sitezttav1.TokenResponse, error) {
	//panic("implement me")
	return nil, nil
}

func (h *handler) ValidateToken(ctx context.Context, req *sitezttav1.TokenRequest) (*sitezttav1.AuthInfo, error) {
	//panic("implement me")
	return nil, nil
}
