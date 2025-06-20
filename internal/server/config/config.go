package config

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Logger LoggerConfig `json:"logger"`
	Server ServerConfig `json:"server"`
	Redis  RedisConfig  `json:"redis"`
}

type LoggerConfig struct {
	Level     slog.Level `json:"level"`
	AddSource bool       `json:"add_source"`
}

type ServerConfig struct {
	Address    string `json:"address"`
	TimeoutS   int64  `json:"timeout_s"`
	TTLS       int64  `json:"ttl_s"`
	Difficulty int    `json:"difficulty"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

func New(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
