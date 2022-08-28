package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type KarmaVote struct {
	DataType       string        `json:"type"`
	VoteId         int           `json:"voteId"`
	Value          int           `json:"voteValue"`
	VotedMessageId int           `json:"votedMessageId"`
	VoteUser       tgbotapi.User `json:"voteUser"`
	VoteTargetUser tgbotapi.User `json:"voteTargetUser"`
	Chat           tgbotapi.Chat `json:"chat"`
}

func voteUp(bot *tgbotapi.BotAPI, message tgbotapi.Message) {

	if !checkVote(bot, message) {
		return
	}

	putVote(message, 1)
	newKarma := getKarma(int(message.ReplyToMessage.From.ID), int(message.Chat.ID))
	responseConfig := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("<b>%s</b> збільшив репутацію <b>%s</b> до <b>%d</b> за <a href='https://t.me/c/%d/%d'>це</a>",
			getName(message.From),
			getName(message.ReplyToMessage.From),
			newKarma,
			-message.Chat.ID-1000000000000,
			message.ReplyToMessage.MessageID,
		))
	responseConfig.ParseMode = "HTML"

	response, _ := bot.Send(responseConfig)
	deleteMessageWithDelay(bot, message, response)
}

func checkVote(bot *tgbotapi.BotAPI, message tgbotapi.Message) bool {

	vote := findVote(message)

	if vote != nil {
		responseConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			fmt.Sprintf("Ви вже голосували за це повідомлення, %s", getName(message.From)),
		)
		bot.Send(responseConfig)
		return false
	}

	if message.From.ID == message.ReplyToMessage.From.ID {
		responseConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			fmt.Sprintf("Не можна за себе голосувати, %s", getName(message.From)),
		)
		bot.Send(responseConfig)
		return false
	}
	//TODO: check multiple votes
	return true
}

func voteDown(bot *tgbotapi.BotAPI, message tgbotapi.Message) {
	if !checkVote(bot, message) {
		return
	}
	putVote(message, -1)
	newKarma := getKarma(int(message.ReplyToMessage.From.ID), int(message.Chat.ID))
	responseConfig := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("<b>%s</b> зменьшив репутацію <b>%s</b> до <b>%d</b> за <a href='https://t.me/c/%d/%d'>це</a>",
			getName(message.From),
			getName(message.ReplyToMessage.From),
			newKarma,
			-message.Chat.ID-1000000000000,
			message.ReplyToMessage.MessageID,
		))
	responseConfig.ParseMode = "HTML"
	response, _ := bot.Send(responseConfig)

	deleteMessageWithDelay(bot, message, response)

}
