package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/xamust/authserver/internal/xlogger"
	"log/slog"
	"time"
)

const (
	confEnvFile = "CONF_ENV_FILE"
)

type Config struct {
	Host           string         `yaml:"host" default:"localhost"`
	Port           int            `yaml:"port" default:"8080"`
	TokenTTL       time.Duration  `yaml:"token_ttl" required:"true"`
	LogLevel       slog.Level     `yaml:"loglevel" default:"debug"`
	LogEnv         xlogger.LogEnv `yaml:"logEnv" env_required:"true"`
	GRPCConfig     GRPCConfig     `yaml:"grpc"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
}

type GRPCConfig struct {
	Host    string `yaml:"host" default:"localhost"`
	Port    int    `yaml:"port" default:"8081"`
	Gateway int    `yaml:"gateway_port" default:"8082"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     int    `yaml:"port" default:"5432"`
	DBName   string `yaml:"db_name" required:"true"`
	User     string `yaml:"user" required:"true"`
	Password string `yaml:"password" required:"true"`
}

func MustLoad() *Config {
	var (
		confFile string
		exist    bool
	)
	if confFile, exist = Env(); !exist {
		confFile = Flag()
	}
	return parseConfig(confFile)
}

func parseConfig(configPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil
	}
	return &cfg
}
