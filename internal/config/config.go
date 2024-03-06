package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultConf = "config/config.yml"

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

// MustLoad loads config into Config struct and returns it
func MustLoad() *Config {
	configPath := filepath.FromSlash(os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		configPath = filepath.FromSlash(DefaultConf)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
