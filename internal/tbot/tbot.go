package tbot

import (
	"context"
	"os"
	"os/signal"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/config"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/storage"
	"go.uber.org/zap"
)

var log *zap.Logger

// Start performs telegram bot setup and runs it
func Start(cfg config.Bot, zlog *zap.Logger, storage storage.Storage) error {
	log = zlog

	var ctxData = newCtxData(storage)

	ctx := context.WithValue(context.Background(), ctxKey, ctxData)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(logMessage),
		bot.WithDefaultHandler(handleDefault),
		bot.WithErrorsHandler(handleErrors),
	}

	b, err := bot.New(cfg.Token, opts...)
	if nil != err {
		return err
	}

	re := regexp.MustCompile(`^/start( [0-9]+)?$`)
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, re, handleCmdStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/poll", bot.MatchTypeExact, handleCmdPoll)

	b.Start(ctx)

	log.Info("Bot stopped")
	return nil
}
