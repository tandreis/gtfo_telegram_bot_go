package main

import (
	dlog "log"

	"github.com/tandreis/gtfo_telegram_bot_go/internal/config"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/logger"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/storage"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/tbot"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		dlog.Fatal(err)
	}

	log := logger.MustInit(cfg.Logger.Level)
	defer log.Sync()

	storage, err := storage.New(cfg.Storage.Type, cfg.Storage.Path)
	if err != nil {
		log.Fatal("Storage init error", zap.Error(err))
	}

	log.Info("Starting GTFO Telegram Bot")

	if err := tbot.Start(cfg.Bot, cfg.Steam, log, storage); err != nil {
		log.Error("Error while running tbot", zap.Error(err))
	}
}
