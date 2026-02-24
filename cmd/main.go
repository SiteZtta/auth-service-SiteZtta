package main

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/internal/app"
	"auth-service-SiteZtta/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		panic(err)
	}
	log := logger.SetupLogger(cfg.Env)
	log.Info("downloaded congig", slog.String("cfgEnv", cfg.Env), slog.Any("cfg", cfg))
	application := app.New(log, cfg, "")
	go application.GRPCServer.MustRun()
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signa := <-stop
	log.Info("shutting down server...", slog.Any("signal", signa.String()))
	application.GRPCServer.Stop()
	log.Info("application stopped")
}
