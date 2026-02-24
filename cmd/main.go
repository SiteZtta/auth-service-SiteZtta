package main

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/pkg/logger"
	"log/slog"
)

func main() {
	cfg, err := config.MustLoad("local")
	if err != nil {
		panic(err)
	}
	log := logger.SetupLogger(cfg.Env)
	log.Info("downloaded congig", slog.String("cfgEnv", cfg.Env), slog.Any("cfg", cfg))
}
