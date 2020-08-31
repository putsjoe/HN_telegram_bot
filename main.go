package main

import (
	"hackernews-telegram/hackernews"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/sayhi"),
		tgbotapi.NewKeyboardButton("/status"),
		tgbotapi.NewKeyboardButton("/latest"),
	),
)

func sayHello(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(10 * time.Second)
	for _ = range ticker.C {
		hello := tgbotapi.NewMessage(361377774, "Hi there mofocker")
		bot.Send(hello)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("1171278568:AAHizulfKvIfASaC0YCKdlXyl0ZFLzRfmwE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	go sayHello(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			case "latest":
				msg.Text = hackernews.PrintItems()
			default:
				msg.Text = "I don't know that command"
				msg.ReplyMarkup = numericKeyboard
			}
			bot.Send(msg)
			continue
		}

	}
}
