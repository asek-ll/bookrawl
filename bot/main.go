package main

import (
	"bookrawl/app/tgbot"
	"bookrawl/scheduler/utils"
	"bookrawl/app/abooks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
)

func main() {
	botToken := os.Getenv("TG_BOT_TOKEN")
	webHookHost := os.Getenv("TG_BOT_WEB_HOOK_HOST")
	webHookPort := os.Getenv("TG_BOT_WEB_HOOK_PORT")

	if webHookPort == "" {
		webHookPort = "443"
	}

	mongoClient, err := utils.CreateMongoClient()
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}


	bot, err := tgbot.CreateBot(botToken)
	if err != nil {
		log.Fatalf("Can't create bot api: %v", err)
	}

	err = setCommands(bot, []tgbotapi.BotCommand{
		tgbotapi.BotCommand{
			Command:     "/list",
			Description: "List of last processed books",
		},
	})
	if err != nil {
		log.Fatalf("Can't set bot commands: %v", err)
	}

	updates, err := listenForWebhook(bot, webHookHost, webHookPort, botToken)

	if err != nil {
		log.Fatal(err)
	}

	bookStore := &abooks.AbookStore{
		Collection: mongoClient.Database("bookrawl").Collection("abooks"),
	}

	for update := range updates {
		log.Printf("%+v\n", update)
		msg := update.Message
		cmd := msg.CommandWithAt()
		if cmd == "list" {
			result, err := bookStore.Find(nil, 20)
			if err != nil {
				log.Fatal(err)
			}

			lines := []string{}

			for _, book := range result.Books {
				lines = append(lines, fmt.Sprintln(book.Author, "-", book.Title, "-", book.Date, book.AuthorId, book.Link))
			}

			sendMessage(bot, msg, strings.Join(lines, "\n"))

		}
	}

}

func listenForWebhook(tgBot *tgbot.TgBot, webHookHost string, webHookPort string, webHookPath string) (tgbotapi.UpdatesChannel, error) {
	content, err := ioutil.ReadFile("./bot-certs/cert.pem")
	if err != nil {
		return nil, err
	}

	cert := tgbotapi.FileBytes{Name: "cert.pem", Bytes: content}

	webHookAddress := "https://" + webHookHost + ":" + webHookPort + "/" + webHookPath
	log.Printf("Use webhook %v", webHookAddress)
	wh, _ := tgbotapi.NewWebhookWithCert(webHookAddress, cert)

	_, err = tgBot.BotApi.Request(wh)
	if err != nil {
		return nil, err
	}

	info, err := tgBot.BotApi.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Callback done")
	})

	updates := tgBot.BotApi.ListenForWebhook("/" + webHookPath)

	go listenAndServeTLS(webHookPort)

	return updates, nil
}

func listenAndServeTLS(webHookPort string) {
	err := http.ListenAndServeTLS(":"+webHookPort, "./bot-certs/cert.pem", "./bot-certs/key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setCommands(tgBot *tgbot.TgBot, commands []tgbotapi.BotCommand) error {
	commandsConfig := tgbotapi.SetMyCommandsConfig{
		Commands:     commands,
		Scope:        nil,
		LanguageCode: "ru",
	}
	resp, err := tgBot.BotApi.Request(commandsConfig)

	if err != nil {
		return err
	}

	log.Printf("Set bot commands %v", string(resp.Result))

	return nil
}

func sendMessage(tgBot *tgbot.TgBot, baseMessage *tgbotapi.Message, text string) error {
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
