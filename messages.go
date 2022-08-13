package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatMessage struct {
	DataType  string           `json:"type"`
	MessageId int              `json:"id"`
	UserId    int64            `json:"userId"`
	Text      string           `json:"text"`
	Date      int64            `json:"date"`
	ReplyToId int64            `json:"replyToId"`
	User      tgbotapi.User    `json:"user"`
	Message   tgbotapi.Message `json:"messageDetails"`
}

var AllowedCommands = []string{
	"/stats",
	"/showRanksTest",
	"/rankDiag",
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	db := getUpdatesDB()
	fmt.Printf("update %+v\n", update)

	docId := fmt.Sprintf("tgupdate:%d", update.UpdateID)

	if update.Message != nil {

		replyToId := 0
		if update.Message.ReplyToMessage != nil {
			replyToId = update.Message.ReplyToMessage.MessageID
		}

		db.Put(context.TODO(), docId, ChatMessage{
			DataType:  "message",
			MessageId: update.Message.MessageID,
			UserId:    int64(update.Message.From.ID),
			Text:      update.Message.Text,
			User:      *update.Message.From,
			ReplyToId: int64(replyToId),
			Date:      int64(update.Message.Date),
			Message:   *update.Message,
		})

		processMessage(bot, *update.Message)
	}

}

func processMessage(bot *tgbotapi.BotAPI, message tgbotapi.Message) {

	messageChunks := strings.Split(message.Text, " ")

	if len(messageChunks) >= 1 && stringInSlice(messageChunks[0], AllowedCommands) {
		switch messageChunks[0] {
		case "/stats":
			stats := getStats(int(message.Chat.ID))
			statsFormatted := "Хуїстика:\n"
			for _, v := range stats {
				statsFormatted += fmt.Sprintf("%s %s: повідомлень: %d, карма: %d\n", calculateDesignation(v.Karma, int(v.MessageCount)), getName(&v.User), v.MessageCount, v.Karma)
			}
			responseConfig := tgbotapi.NewMessage(
				message.Chat.ID,
				statsFormatted,
			)
			// TODO: Investigate escaping options
			//responseConfig.ParseMode = "MarkdownV2"
			bot.Send(responseConfig)

		case "/showRanksTest":
			for _, chunk := range getRanksDescriptionsSplitted() {

				rankDesk := ""
				for _, v := range chunk {
					rankDesk = rankDesk + fmt.Sprintf("%s -> Rank: %s - minmsg: %d congrat:%s  pic: %s\n",
						v.Emoji,
						v.RankName,
						v.MinMessages,
						v.CongratulationMessage,
						v.Picture,
					)
				}

				responseConfig := tgbotapi.NewMessage(
					message.Chat.ID,
					rankDesk,
				)
				bot.Send(responseConfig)
			}
			for _, chunk := range getKarmaRanksDescriptionsSplitted() {

				rankDesk := ""
				for _, v := range chunk {
					rankDesk = rankDesk + fmt.Sprintf("%d -> %s\n",
						v.Karma,
						v.Rank,
					)
				}

				responseConfig := tgbotapi.NewMessage(
					message.Chat.ID,
					rankDesk,
				)
				bot.Send(responseConfig)
			}
		case "/rankDiag":
			if len(messageChunks) != 3 {
				return
			} else {
				s, err := strconv.Atoi(messageChunks[1])
				if err != nil {
					log.Printf("Incorrect Params %s", message.Text)
				}
				k, err := strconv.Atoi(messageChunks[2])
				if err != nil {
					log.Printf("Incorrect Params %s", message.Text)
				}
				responseConfig := tgbotapi.NewMessage(
					message.Chat.ID,
					calculateDesignation(k, s),
				)
				bot.Send(responseConfig)

			}
		}
	}

	if stringInSlice(message.Text, []string{"+", "-"}) && message.ReplyToMessage != nil && !message.ReplyToMessage.From.IsBot {
		switch message.Text {
		case "+":
			voteUp(bot, message)
		case "-":
			voteDown(bot, message)
		}

	}
}

func deleteMessageWithDelay(bot *tgbotapi.BotAPI, message tgbotapi.Message, response tgbotapi.Message) {

	go func(message tgbotapi.Message, response tgbotapi.Message) {
		time.Sleep(10 * time.Second)
		delete := tgbotapi.NewDeleteMessage(message.Chat.ID, response.MessageID)
		bot.Send(delete)
		delete = tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		bot.Send(delete)
	}(message, response)

}

func getName(user *tgbotapi.User) string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
