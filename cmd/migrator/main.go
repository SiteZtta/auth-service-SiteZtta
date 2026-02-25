package main

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/pkg/logger"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // init postgres driver
)

const (
	defaultDb = "postgres"
)

func main() {
	var cfgPath, migrationsPath, migrationsTable string
	flag.StringVar(&cfgPath, "cfg", "", "path to cfg dir")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of migrations table")
	flag.Parse()
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		panic(err)
	}
	log := logger.SetupLogger(cfg.Env)
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	ensureDbCreated(cfg, log)
	m, err := migrate.New(
		"file://"+migrationsPath,
		config.GetConnStringMigrate(cfg)+"&x-migrations-table="+migrationsTable+"",
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply")
		} else {
			panic(err)
		}
	} else {
		log.Info("migrations applied successfully")
	}
}

func ensureDbCreated(cfg config.Config, log *slog.Logger) error {
	cfgDef := cfg // config for default db
	cfgDef.Database.Name = defaultDb
	connStr := config.GetConnString(cfgDef)
	dbAdmin, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	var isExists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
	if err = dbAdmin.QueryRow(query, cfg.Database.Name).Scan(&isExists); err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}
	if !isExists {
		_, err = dbAdmin.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.Database.Name))
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	} else {
		log.Info("Database already exists", "name", cfg.Database.Name)
	}
	return nil
}
