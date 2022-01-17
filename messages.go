package main

import (
	"context"
	_ "embed"
	"fmt"
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

	if stringInSlice(message.Text, AllowedCommands) {
		switch message.Text {
		case "/stats":
			stats := getStats(int(message.Chat.ID))
			statsFormatted := "Статистика-Хуистика:\n"
			for _, v := range stats {
				statsFormatted += fmt.Sprintf("%s: cообщений: %d, карма: %d\n", getName(&v.User), v.MessageCount, v.Karma)
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
			// TODO: Investigate escaping options
			//responseConfig.ParseMode = "MarkdownV2"

			//TODO: User rank from db
			// case "/importRanks":
			// 	metaDb := getUserMetadataDB()
			// 	getRanksCSVEmbeded()

			//case "/rankDiag":
			//errors := ""
			//docs := ranksDb.AllDocs(context.TODO())

			//log.Printf("%+v", docs)

			// 	dir, err := os.MkdirTemp("", "srakabot")
			// 	if err != nil {
			// 		return
			// 	}

			// 	for i, v := range getRanksCSVEmbeded() {
			// 		res, err := http.DefaultClient.Get(v.Picture)
			// 		log.Printf("%+v", res)
			// 		if err != nil || res.StatusCode >= 400 {
			// 			errors = errors + fmt.Sprintf("Ранг %s: кортинко %s не грузиццо\n", v.RankName, v.Picture)
			// 		} else {

			// 			b, err := ioutil.ReadAll(res.Body)
			// 			if err != nil {
			// 				continue
			// 			}
			// 			ct := strings.Split(http.DetectContentType(b), "/")

			// 			ioutil.WriteFile(fmt.Sprintf("%s/rank_%d.%s", dir, i, ct[1]), b, 0644)

			// 			//	tgbotapi.RequestFileData{}
			// 			//	tgbotapi.NewPhoto(message.Chat.ID, res.Body)
			// 		}
			// 		if i%15 == 0 {
			// 			if errors != "" {
			// 				responseConfig := tgbotapi.NewMessage(
			// 					message.Chat.ID,
			// 					errors,
			// 				)
			// 				bot.Send(responseConfig)
			// 				errors = ""
			// 			}
			// 		}
			// 	}
			// 	if errors != "" {
			// 		responseConfig := tgbotapi.NewMessage(
			// 			message.Chat.ID,
			// 			errors,
			// 		)
			// 		bot.Send(responseConfig)
			// 		errors = ""
			// 	}
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
