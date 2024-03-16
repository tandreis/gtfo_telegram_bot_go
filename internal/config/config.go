package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger
	Storage
	Bot
	Steam
}

type Logger struct {
	Level string `yaml:"level" env-default:"info"`
}

type Storage struct {
	Type string `yaml:"type" env-required:"true"`
	Path string `yaml:"path"`
}

type Bot struct {
	Token    string `yaml:"token" env:"BOT_TOKEN" env-required:"true"`
	Args     string `yaml:"args"`
	MaxPolls int    `yaml:"max_polls" env-default:"1"`
}

type User struct {
	Name       string `yaml:"name" env-required:"true"`
	SteamID    string `yaml:"steam_id" env-required:"true"`
	TelegramID int64  `yaml:"telegram_id" env-required:"true"`
	ChatID     int64  `yaml:"chat_id" env-required:"true"`
}

type Steam struct {
	ApiKey string `yaml:"api_key" env:"STEAM_TOKEN"`
	Users  []User `yaml:"users"`
}

// Load read config from file and returns Config struct.
// On error returns nil and error
func Load() (*Config, error) {
	configPath := filepath.FromSlash(os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("config file does not exists: %s", configPath))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.New(fmt.Sprintf("cannot read config: %s", err))
	}

	return &cfg, nil
}
