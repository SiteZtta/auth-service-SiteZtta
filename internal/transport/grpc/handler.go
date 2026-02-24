package grpc

import (
	"auth-service-SiteZtta/internal/transport/grpc/v1/dto"
	"context"

	sitezttav1 "github.com/SiteZtta/protos-SiteZtta/gen/go/auth"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	CreateUser(ctx context.Context, in dto.SignUpInput) (int64, error)
	GenerateToken(ctx context.Context, in dto.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (dto.AuthInfo, error)
}

type handler struct {
	sitezttav1.UnimplementedAuthServiceServer
	validator *validator.Validate
}

func Register(gRPC *grpc.Server) {
	sitezttav1.RegisterAuthServiceServer(gRPC, &handler{validator: validator.New()})
}

func (h *handler) CreateUser(ctx context.Context, req *sitezttav1.SignUpRequest) (*sitezttav1.UserIdResponse, error) {
	// validation
	input := dto.SignUpInput{
		UserName: req.GetUserName(),
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//panic("implement me")
	return nil, nil
}

func (h *handler) GenerateToken(ctx context.Context, req *sitezttav1.SignInRequest) (*sitezttav1.TokenResponse, error) {
	// validation
	input := dto.SignInInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//panic("implement me")
	return &sitezttav1.TokenResponse{Token: req.GetLogin() + req.GetPassword()}, nil
}

func (h *handler) ValidateToken(ctx context.Context, req *sitezttav1.TokenRequest) (*sitezttav1.AuthInfo, error) {
	// validation
	input := dto.TokenInput{
		Token: req.GetToken(),
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//panic("implement me")
	return nil, nil
}
