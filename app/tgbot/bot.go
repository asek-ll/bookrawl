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

func (tgBot *TgBot) ReplyToMessage(baseMessage *tgbotapi.Message, text string) error {
	config := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: baseMessage.Chat.ID,
			ChannelUsername: baseMessage.Chat.UserName,
			ReplyToMessageID: baseMessage.MessageID,
			ReplyMarkup: nil,
			DisableNotification: false,
			AllowSendingWithoutReply: true,
		},
		Text: text,
		ParseMode: "",
		Entities: []tgbotapi.MessageEntity{},
		DisableWebPagePreview: false,
	}
	_, err := tgBot.BotApi.Request(config)

	if err != nil {
		return err
	}

	return nil
}
