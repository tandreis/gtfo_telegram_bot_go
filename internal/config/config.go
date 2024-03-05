package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultConf = "config/config.yml"

type Config struct {
	Bot
	Steam
}

type Bot struct {
	Token    string `yaml:"token" env:"BOT_TOKEN" env-required:"true"`
	Args     string `yaml:"args"`
	MaxPolls int    `yaml:"max_polls" env-default:"1"`
}

type User struct {
	Name       string `yaml:"name" env-required:"true"`
	SteamID    string `yaml:"sid" env-required:"true"`
	TelegramID int    `yaml:"tid" env-required:"true"`
}

type Steam struct {
	ApiKey string `yaml:"api_key" env:"STEAM_TOKEN"`
	Users  []User `yaml:"users"`
}

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
