package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type JWTConfig struct {
	Secret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.AutomaticEnv()

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("DATABASE_DSN"),
		},
		Redis: RedisConfig{
			Addr: viper.GetString("REDIS_ADDR"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
	}, nil
}
