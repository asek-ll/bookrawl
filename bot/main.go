package main

import (
	"bookrawl/app/tgbot"
	"bookrawl/scheduler/utils"
	"bookrawl/app/dao"
	"bookrawl/bot/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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


	daoHolder := dao.NewDaoHolder(mongoClient)

	cmds := []commands.Command{
		&commands.ListCommand{DaoHolder: daoHolder},
		&commands.SearchByAuthorCommand{DaoHolder: daoHolder},
	}
	
	err = setCommands(bot, cmds)

	if err != nil {
		log.Fatalf("Can't set bot commands: %v", err)
	}

	updates, err := listenForWebhook(bot, webHookHost, webHookPort, botToken)

	if err != nil {
		log.Fatal(err)
	}

	processCommands(bot, cmds, updates)
}

func processCommands(tgBot *tgbot.TgBot, cmds []commands.Command, updates tgbotapi.UpdatesChannel) error {

	log.Println("Process commands")

	for update := range updates {
		msg := update.Message
		if msg.IsCommand() {
			cmdName := msg.CommandWithAt()

			log.Println("Get cmd", cmdName)
			ctx := &commands.Context{
				TgBot: tgBot,
				Message: msg,
			}
			for _, cmd := range cmds {
				if cmd.GetName() == cmdName {
					log.Println("Found processor", cmdName)
					err := cmd.Run(ctx)
					if err != nil {
						log.Println(err)
					}
				}

			}
		}
	}

	return nil
}

func listenForWebhook(tgBot *tgbot.TgBot, webHookHost string, webHookPort string, webHookPath string) (tgbotapi.UpdatesChannel, error) {
	content, err := ioutil.ReadFile("/bot-certs/cert.pem")
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

	log.Println("Start listen")
	go listenAndServeTLS(webHookPort)
	log.Println("Return updates")

	return updates, nil
}

func listenAndServeTLS(webHookPort string) {
	err := http.ListenAndServeTLS(":"+webHookPort, "/bot-certs/cert.pem", "/bot-certs/key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setCommands(tgBot *tgbot.TgBot, commands []commands.Command) error {

	botCommands := make([]tgbotapi.BotCommand, len(commands))

	for i, cmd := range commands {
		botCommands[i] = tgbotapi.BotCommand{
			Command:     "/" + cmd.GetName(),
			Description: cmd.GetDescription(),
		}
	}


	commandsConfig := tgbotapi.SetMyCommandsConfig{
		Commands:     botCommands,
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
