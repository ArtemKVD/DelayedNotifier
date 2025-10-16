package config

import (
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

type Config struct {
	HTTPPort      string
	RabbitMQURL   string
	RedisURL      string
	TelegramToken string
	PostgreSQL    PostgreSQLConfig
}

type PostgreSQLConfig struct {
	MasterDSN string
	SlaveDSNs []string
}

func Load() *Config {
	cfg := config.New()
	err := cfg.Load("config.yaml", ".env", "NOTIFICATION")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Load config error")
	}

	var appConfig Config
	err = cfg.Unmarshal(&appConfig)
	if err != nil {
		return &Config{
			HTTPPort:      cfg.GetString("http_port"),
			RabbitMQURL:   cfg.GetString("rabbitmq_url"),
			RedisURL:      cfg.GetString("redis_url"),
			TelegramToken: cfg.GetString("telegram_token"),
			PostgreSQL: PostgreSQLConfig{
				MasterDSN: cfg.GetString("postgresql.master_dsn"),
				SlaveDSNs: cfg.GetStringSlice("postgresql.slave_dsns"),
			},
		}
	}

	return &appConfig
}
