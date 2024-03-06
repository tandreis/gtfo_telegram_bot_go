package tbot

import (
	"fmt"
	"html"

	"github.com/go-telegram/bot/models"
)

func mentionUser(user *models.User) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`,
		user.ID, html.EscapeString(user.FirstName))
}
