package main

import (
	"hackernews-telegram/hackernews"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/sayhi"),
		tgbotapi.NewKeyboardButton("/status"),
		tgbotapi.NewKeyboardButton("/latest"),
	),
)

// Shall be used to send the latest unread items every hour or so.
func sayHello(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		hello := tgbotapi.NewMessage(361377774, "Hi there mofocker")
		bot.Send(hello)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("1171278568:AAHizulfKvIfASaC0YCKdlXyl0ZFLzRfmwE")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go hackernews.UpdatePosts()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID != 361377774 {
			continue
		}

		log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			if strings.Contains(update.Message.Command(), "add_") {
				msg.Text = strings.Replace(
					update.Message.Command(), "add_", "", 1)
				bot.Send(msg)
				continue
			}

			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."

			case "ping":
				msg.Text = "pong"

			case "latest":
				msg.Text = hackernews.TextItems()

			default:
				msg.Text = "I don't know that command"
				msg.ReplyMarkup = numericKeyboard
			}
			bot.Send(msg)
			continue
		}

	}
}
