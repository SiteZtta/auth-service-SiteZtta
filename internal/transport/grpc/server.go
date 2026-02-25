package grpc

import (
	"auth-service-SiteZtta/internal/service/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *Server {
	authService := auth.New(log)
	gRPCServer := grpc.NewServer()
	Register(gRPCServer)
	return &Server{log: log, gRPCServer: gRPCServer, port: port}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const fn = "auth-service-SiteZtta.internal.transport.grpc.run"
	log := s.log.With(slog.String("fn", fn), slog.Int("port", s.port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	log.Info("gRPC server is running", slog.String("addr", lis.Addr().String()))
	if err = s.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Server) Stop() {
	const fn = "auth-service-SiteZtta.internal.transport.grpc.stop"
	s.log.With(slog.String("fn", fn)).Info("Stopping gRPC server")
	s.gRPCServer.GracefulStop()
}
