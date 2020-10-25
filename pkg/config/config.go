package config

import (
	"os"
)

type RedisConfig struct {
	Addr      string
	KeyPrefix string
	Network   string
}

type TwitchConfig struct {
	ClientID     string
	ClientSecret string
}

type Config struct {
	ServerAddr string
	Redis      *RedisConfig
	Twitch     *TwitchConfig
}

func New() *Config {
	return &Config{
		ServerAddr: optionalEnv("SERVER_ADDR", ":4500"),

		Redis: &RedisConfig{
			Addr:      optionalEnv("REDIS_ADDR", "127.0.0.1:6379"),
			KeyPrefix: optionalEnv("REDIS_KEY_PREFIX", "customapi::"),
			Network:   optionalEnv("REDIS_NETWORK", "tcp"),
		},

		Twitch: &TwitchConfig{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}

func optionalEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
