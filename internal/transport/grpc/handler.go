package grpc

import (
	"auth-service-SiteZtta/internal/service/auth"
	"auth-service-SiteZtta/internal/storage"
	"auth-service-SiteZtta/internal/transport/grpc/v1/dto"
	"context"
	"errors"

	sitezttav1 "github.com/SiteZtta/protos-SiteZtta/gen/go/auth"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	CreateUser(ctx context.Context, in dto.SignUpInput) (uid int64, err error)
	GenerateToken(ctx context.Context, in dto.SignInInput) (token string, err error)
	ValidateToken(ctx context.Context, token string) (authInfo dto.AuthInfo, err error)
}

type handler struct {
	sitezttav1.UnimplementedAuthServiceServer
	validator *validator.Validate
	auth      Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sitezttav1.RegisterAuthServiceServer(gRPC, &handler{validator: validator.New(), auth: auth})
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
	// Business logic
	userId, err := h.auth.CreateUser(ctx, input)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sitezttav1.UserIdResponse{UserId: userId}, nil
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
	// Business logic
	token, err := h.auth.GenerateToken(ctx, input)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sitezttav1.TokenResponse{Token: token}, nil
}

func (h *handler) ValidateToken(ctx context.Context, req *sitezttav1.TokenRequest) (*sitezttav1.AuthInfo, error) {
	// validation
	input := dto.TokenInput{
		Token: req.GetToken(),
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Business logic
	authInfo, err := h.auth.ValidateToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	prtotobufAuthInfo := &sitezttav1.AuthInfo{
		UserId: authInfo.UserId,
		Role:   sitezttav1.Role(authInfo.Role),
	}
	return prtotobufAuthInfo, nil
}
