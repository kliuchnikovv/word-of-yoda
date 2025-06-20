package config

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Logger LoggerConfig `json:"logger"`
	Server string       `json:"server"`
	Client ClientConfig `json:"client"`
}

type LoggerConfig struct {
	Level     slog.Level `json:"level"`
	AddSource bool       `json:"add_source"`
}

type ClientConfig struct {
	TimeoutS int64 `json:"timeout_s"`
	SleepS   int64 `json:"sleep_s"`
}

func New(path string) (*Config, error) {
	file, err := os.Open("/app/config/client_config_test.json")
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
