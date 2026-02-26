package app

import (
	"auth-service-SiteZtta/config"
	"log/slog"

	"auth-service-SiteZtta/internal/service/auth"
	"auth-service-SiteZtta/internal/storage/pgdb"
	transGrpc "auth-service-SiteZtta/internal/transport/grpc"
)

type App struct {
	GRPCServer *transGrpc.Server
}

func New(log *slog.Logger, cfg config.Config, connStr string) *App {
	storage, err := pgdb.New(connStr)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, cfg.Auth)
	grpcServer := transGrpc.New(log, cfg.GrpcServer.Port, authService)
	return &App{GRPCServer: grpcServer}
}
