package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	//connect to touchbase
	cbConnect()

	token := os.Getenv("TGBOT_API_TOKEN")
	file_location := os.Getenv("HOME") + "/.srakabot_token"
	if token == "" {
		contents, err := os.ReadFile(file_location)
		if err == nil {
			token = strings.Trim(string(contents), "\n")
		} else {
			log.Printf("%s could not be read and environment variable TG_BOT_API_TOKEN is not defined", file_location)
		}
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		os.Exit(1)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		processUpdate(bot, update)
	}
}
