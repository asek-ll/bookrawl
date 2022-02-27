package tgbot


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type TgBot struct {
	BotApi *tgbotapi.BotAPI
}


func CreateBot(token string) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TgBot{BotApi: bot}, nil
}
