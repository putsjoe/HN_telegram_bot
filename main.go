package main

import (
	"fmt"
	"hackernews-telegram/hackernews"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/list"),
		tgbotapi.NewKeyboardButton("/latest"),
	),
)

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
				p := strings.Replace(
					update.Message.Command(), "add_", "", 1)
				postID, err := strconv.Atoi(p)
				if err != nil {
					fmt.Println(err)
				}
				hackernews.SavePost(int(update.Message.Chat.ID), postID)
				msg.Text = "Saved " + p
				bot.Send(msg)
				continue
			}

			if strings.Contains(update.Message.Command(), "del_") {
				p := strings.Replace(
					update.Message.Command(), "del_", "", 1)
				postID, err := strconv.Atoi(p)
				if err != nil {
					fmt.Println(err)
				}
				hackernews.DeletePost(int(update.Message.Chat.ID), postID)
				msg.Text = "Deleted " + p
				bot.Send(msg)
				continue
			}

			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."

			case "ping":
				u, t := hackernews.UnreadItems()
				msg.Text = fmt.Sprintf("pong - Unread: %d  Total: %d", u, t)

			case "latest":
				msg.Text = hackernews.TextItems()

			case "saves":
			case "list":
				msg.Text = hackernews.GetSavedPosts(int(update.Message.Chat.ID))

			default:
				msg.Text = "I don't know that command"
				msg.ReplyMarkup = numericKeyboard
			}
			bot.Send(msg)
			continue
		}

	}
}
