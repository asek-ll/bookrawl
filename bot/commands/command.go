package commands

import (
	"bookrawl/app/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	TgBot *tgbot.TgBot
	Message *tgbotapi.Message
}

func (ctx *Context) Reply(msg string) error {
	return ctx.TgBot.ReplyToMessage(ctx.Message, msg)
}


type Command interface {
	GetName() string
	GetDescription() string
	Run(*Context) error
}
