package main

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/internal/app"
	"auth-service-SiteZtta/pkg/logger"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cfgPath string
	// --cfg="./config/local.yaml"
	flag.StringVar(&cfgPath, "cfg", "", "path to cfg dir")
	flag.Parse()
	if cfgPath == "" {
		panic("cfg path is empty")
	}
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		panic(err)
	}
	log := logger.SetupLogger(cfg.Env)
	log.Info("downloaded congig", slog.String("cfgEnv", cfg.Env), slog.Any("cfg", cfg))
	connStr := config.GetConnString(cfg)
	application := app.New(log, cfg, connStr)
	go application.GRPCServer.MustRun()
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signa := <-stop
	log.Info("shutting down server...", slog.Any("signal", signa.String()))
	application.GRPCServer.Stop()
	log.Info("application stopped")
}
