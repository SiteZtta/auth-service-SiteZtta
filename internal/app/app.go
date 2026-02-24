package app

import (
	"auth-service-SiteZtta/config"
	"log/slog"

	transGrpc "auth-service-SiteZtta/internal/transport/grpc"
)

type App struct {
	GRPCServer *transGrpc.Server
}

func New(log *slog.Logger, cfg config.Config, connStr string) *App {
	// TODO: init storage
	// TODO: init auth service (auth)
	grpcServer := transGrpc.New(log, cfg.GrpcServer.Port)
	return &App{GRPCServer: grpcServer}
}
