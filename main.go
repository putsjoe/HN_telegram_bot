package main

import (
	"flag"
	"fmt"
	"hackernews-telegram/hackernews"
	"log"
	"os"
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

func newMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: false,
	}
}

func main() {
	apiKey := flag.String("key", "", "The Telegram bot API key")
	userID := flag.Int("user", 0, "Restrict bot to given user ID")
	flag.Parse()

	if *apiKey == "" {
		fmt.Println("No API key supplied, use the -key flag")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(*apiKey)
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
		if *userID > 0 && update.Message.Chat.ID != int64(*userID) {
			continue
		}

		log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

		if update.Message.Text == "L" {
			msg := newMessage(update.Message.Chat.ID, "")
			msg.Text = hackernews.GetSavedPosts(int(update.Message.Chat.ID), true)
			bot.Send(msg)
			continue
		}
		if update.Message.Text == "P" {
			msg := newMessage(update.Message.Chat.ID, "")
			u, t, s := hackernews.UnreadItems()
			msg.Text = fmt.Sprintf("Unread: %d  Saved: %d  Total: %d", u, s, t)
			bot.Send(msg)
			continue
		}

		if update.Message.IsCommand() {
			msg := newMessage(update.Message.Chat.ID, "")

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
				u, t, s := hackernews.UnreadItems()
				msg.Text = fmt.Sprintf("Unread: %d  Saved: %d  Total: %d", u, s, t)

			case "latest":
				msg.Text = hackernews.TextItems()

			case "saves":
			case "list":
				msg.Text = hackernews.GetSavedPosts(int(update.Message.Chat.ID), false)
			default:
				msg.Text = "I don't know that command"
				msg.ReplyMarkup = numericKeyboard
			}
			if _, err := bot.Send(msg); err != nil {
				log.Println("bot.Send: ", err)
				log.Println("Message: ", msg)
			}

			continue
		}

	}
}
