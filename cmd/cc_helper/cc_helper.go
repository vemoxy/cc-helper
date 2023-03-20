package main

import (
	"log"

	"github.com/vemoxy/cc-helper/config"
	"github.com/vemoxy/cc-helper/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Telegram.ApiKey)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Print(err.Error())
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		var replyMessageText string

		command := update.Message.Command()
		switch command {
		case "checkmcc":
			if len(update.Message.CommandArguments()) == 0 {
				replyMessageText = "Please input a shop name."
				break
			}
			results := handler.CheckMcc(update.Message.CommandArguments())
			replyMessageText = handler.GenerateMccV3ResultMessage(results)
		default:
			replyMessageText = "Please use an existing command."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessageText)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
