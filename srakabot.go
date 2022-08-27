package main

import (
	"flag"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var appConfig *Config

var AllowedCommands = []string{
	"/stats",
	"/showRanksTest",
	"/rankDiag",
	"/bledina",
}

func main() {

	var err error
	configFileLocationPtr := flag.String("config", os.Getenv("HOME")+"/.srakabot.yml", "Config file location")

	appConfig, err = getConfig(*configFileLocationPtr)
	if err != nil {
		log.Fatalf("Could not load configuration")
	}

	//connect to touchbase
	cbConnect()
	if !sraketa_init() {
		log.Printf("Could not initialize alerts")
	}

	bot, err := tgbotapi.NewBotAPI(appConfig.Telegram.ApiToken)
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
