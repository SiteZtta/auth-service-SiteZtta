package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `mapstructure:"env"`
	Database   Database   `mapstructure:"database"`
	GrpcServer GrpcServer `mapstructure:"grpc_server"`
	Auth       Auth       `mapstructure:"auth"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type GrpcServer struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Auth struct {
	SigningKey string        `mapstructure:"signing_key"`
	TokenTtl   time.Duration `mapstructure:"token_ttl"`
}

func MustLoad() (Config, error) {
	path := fetchCfgDirPath()
	if path == "" {
		return Config{}, errors.New("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file not found on path %s: %w", path, err)
	}
	// Setting viper
	viper.SetConfigFile(path)
	// Env variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // ex: http_server.port -> APP_HTTP_SERVER_PORT
	// Reading config
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("fatal error reading config file: %w", err)
	}
	if err := validateRequired(); err != nil {
		return Config{}, fmt.Errorf("Error validating config: %w", err)
	}
	cfg := Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("fatal error unmarshaling config: %w", err)
	}
	return cfg, nil
}

func fetchCfgDirPath() string {
	var path string
	// --cfg="./config/local.yaml"
	flag.StringVar(&path, "cfg", "", "path to cfg dir")
	flag.Parse()
	if path == "" {
		panic("cfg path is empty")
	}
	return path
}

func validateRequired() error {
	required := []string{
		"env",
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.name",
	}
	for _, field := range required {
		if !viper.IsSet(field) {
			return fmt.Errorf("required field '%s' is not set", field)
		}
	}
	return nil
}

func getConnString(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
}
