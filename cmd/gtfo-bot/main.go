package main

import (
	"github.com/tandreis/gtfo_telegram_bot_go/internal/config"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustInit(cfg.Logger.Level)
	defer log.Sync()

	log.Info("Starting GTFO Telegram Bot")
}
