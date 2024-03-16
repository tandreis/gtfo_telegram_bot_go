package tbot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/tandreis/gtfo_telegram_bot_go/internal/steam"
	strg "github.com/tandreis/gtfo_telegram_bot_go/internal/storage"
	"go.uber.org/zap"
)

var (
	pollQuestion = "Сегодня"
	pollOptions  = []string{"Играем!", "Я не буду", "Не знаю"}
	pollAnswers  = [...]string{"играет!", "не будет играть \u2639\uFE0F", "еще не определился"}
)

func logMessage(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Debug("User message",
				zap.Int64("chat_id", update.Message.Chat.ID),
				zap.Int64("user_id", update.Message.From.ID),
				zap.String("text", update.Message.Text))
		}
		next(ctx, b, update)
	}
}

func fmtDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	months := days / 30

	if months > 0 {
		return fmt.Sprintf("%d мес", months)
	}

	weeks := days / 7
	if weeks > 0 {
		return fmt.Sprintf("%d нед", weeks)
	}
	if days > 0 {
		return fmt.Sprintf("%d д", days)
	}

	hours := int(duration.Hours()) % 24
	if hours > 0 {
		return fmt.Sprintf("%d ч", hours)
	}

	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%d мин", minutes)
}

func handleCmdStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		log.Warn("Got null update message, skipping")
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Unsupported",
	})
}

func handleCmdSteamStatus(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		log.Warn("Got null update message, skipping")
		return
	}

	token, err := getSteamToken(ctx)
	if err != nil {
		log.Error("Failed to get context", zap.Error(err))
		return
	}

	storage, err := getStorage(ctx)
	if err != nil {
		log.Error("Failed to get context", zap.Error(err))
		return
	}

	users, _ := storage.GetUsers(update.Message.Chat.ID)

	var steamIDs []string
	for _, u := range users {
		steamIDs = append(steamIDs, u.SteamID)
	}

	steamUsers := steam.GetPlayerSummaries(token, steamIDs)

	var message = "Статус игроков."
	for _, u := range steamUsers.Response.Players {
		if u.Online {
			message += fmt.Sprintf("\n&#9989;<i>%s</i> сейчас <b>%s</b>.",
				u.Name, u.StateStr)
			if u.GameName != "" {
				message += fmt.Sprintf(" Играет в <a href=\"%s\">%s</a>.",
					u.GameURL, u.GameName)
			}
		} else {
			message += fmt.Sprintf("\n&#10060;<i>%s</i> сейчас <b>%s</b>. (%s)",
				u.Name, u.StateStr, fmtDuration(time.Since(time.Unix(u.LastLogoff, 0))))
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
	})

}

func handleCmdPoll(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		log.Warn("Got null update message, skipping")
		return
	}

	memberCount, err := b.GetChatMemberCount(
		ctx, &bot.GetChatMemberCountParams{ChatID: update.Message.Chat.ID})
	if err != nil {
		log.Error("Member count error", zap.Error(err))
		return
	}

	p := &bot.SendPollParams{
		ChatID:                update.Message.Chat.ID,
		Question:              pollQuestion,
		Options:               pollOptions,
		IsAnonymous:           bot.False(),
		AllowsMultipleAnswers: false,
	}

	message, err := b.SendPoll(ctx, p)
	if err != nil {
		log.Error("Poll send error", zap.Error(err))
		return
	}

	log.Info("Poll started",
		zap.Int64("user_id", update.Message.From.ID),
		zap.Int64("chat_id", update.Message.Chat.ID))

	storage, err := getStorage(ctx)
	if err != nil {
		log.Error("Failed to get context", zap.Error(err))
		return
	}
	storage.CreatePoll(message.Poll.ID, strg.PollEntity{
		ChatID:     update.Message.Chat.ID,
		MaxAnswers: memberCount - 1,
	})
}

func handlePollAnswer(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.PollAnswer.User == nil {
		log.Error("Anonymous polls is not supported")
		return
	}

	var message string
	var answer int = -1

	if len(update.PollAnswer.OptionIDs) != 0 {
		answer = update.PollAnswer.OptionIDs[0]
		if answer < len(pollAnswers) {
			message = fmt.Sprintf("%s %s", mentionUser(update.PollAnswer.User), pollAnswers[answer])
		}
	} else {
		message = fmt.Sprintf("%s передумал голосовать!", mentionUser(update.PollAnswer.User))
	}

	log.Info("Poll answer",
		zap.String("poll_id", update.PollAnswer.PollID),
		zap.Int64("user_id", update.PollAnswer.User.ID),
		zap.Int("answer", answer))

	storage, err := getStorage(ctx)
	if err != nil {
		log.Error("Failed to get context", zap.Error(err))
		return
	}
	poll, _ := storage.GetPoll(update.PollAnswer.PollID)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    poll.ChatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
}

func handlePoll(ctx context.Context, b *bot.Bot, update *models.Update) {
	storage, err := getStorage(ctx)
	if err != nil {
		log.Error("Failed to get context", zap.Error(err))
		return
	}
	poll, _ := storage.GetPoll(update.Poll.ID)

	log.Info("Poll",
		zap.String("poll_id", update.Poll.ID),
		zap.Int64("chat_id", poll.ChatID),
		zap.Int("total_voted", update.Poll.TotalVoterCount))

	if update.Poll.TotalVoterCount >= poll.MaxAnswers {
		b.StopPoll(ctx, &bot.StopPollParams{ChatID: poll.ChatID})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: poll.ChatID,
			Text:   "Все проголосовали.",
		})
		_ = storage.DeletePoll(update.Poll.ID)
	}
}

func handleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.PollAnswer != nil {
		handlePollAnswer(ctx, b, update)
	}

	if update.Poll != nil {
		handlePoll(ctx, b, update)
	}
}

func handleErrors(err error) {
	log.Error("Error in tgbot", zap.Error(err))
}
