package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-kivik/couchdb/v4"
	"github.com/go-kivik/kivik/v4"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var cbInstance *kivik.Client

func cbConnect() {
	client, err := kivik.New("couch", "http://localhost:5984/")
	if err != nil {
		panic("Could not connect to database")
	}
	err = client.Authenticate(context.TODO(), couchdb.BasicAuth("admin", "admin"))
	if err != nil {
		panic(err)
	}
	cbInstance = client
}

func getUpdatesDB() *kivik.DB {
	return cbInstance.DB("srakabot_db")
}

// func getUserMetadataDB() *kivik.DB {
// 	db := cbInstance.DB("srakabot_user_metadata")
// 	return db
// }

func putVote(message tgbotapi.Message, value int) {
	db := getUpdatesDB()
	voteId := fmt.Sprintf("vote:%d", message.MessageID)

	db.Put(context.TODO(), voteId, KarmaVote{
		DataType:       "vote",
		Value:          value,
		VotedMessageId: message.ReplyToMessage.MessageID,
		VoteUser:       *message.From,
		VoteTargetUser: *message.ReplyToMessage.From,
		VoteId:         message.MessageID,
		Chat:           *message.Chat,
	})

}

type KarmaUserKey struct {
	ChatId int `json:"chat"`
	UserId int `json:"user"`
}

func getKarma(userId int, chatId int) int {
	db := getUpdatesDB()
	resultSet := db.Query(context.TODO(), "_design/aggregateByPostCount", "_view/votesAndChatStats", kivik.Options{
		"reduce":   true,
		"group":    true,
		"startkey": []int{chatId, userId},
		"endkey":   []int{chatId, userId},
	})

	if resultSet.Err() != nil {
		fmt.Printf("resultSet problems %v", resultSet.Err())
		return 0
	}
	var value StatsEntry

	for resultSet.Next() {
		var key interface{}
		err := resultSet.ScanValue(&value)

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		fmt.Printf("%+v %+v\n", key, value)
	}

	return value.Karma
}

type StatsEntry struct {
	Karma        int           `json:"karma"`
	MessageCount int64         `json:"messageCount"`
	User         tgbotapi.User `json:"user"`
}

//TODO: filter on view layer
func getStats(chatId int) []StatsEntry {
	db := getUpdatesDB()
	resultSet := db.Query(context.TODO(), "_design/aggregateByPostCount", "_view/votesAndChatStats", kivik.Options{
		"reduce": true,
		"group":  true,
	})
	var entries []StatsEntry

	for resultSet.Next() {
		entry := StatsEntry{}
		key := []int64{0, 0}
		resultSet.ScanKey(&key)
		err := resultSet.ScanValue(&entry)
		if err != nil {
			fmt.Printf("%+v", err)
		}
		if key[0] == int64(chatId) {
			entries = append(entries, entry)
		}
	}
	sort.SliceStable(entries, func(i int, j int) bool {
		return entries[i].MessageCount > entries[j].MessageCount
	})
	return entries
}
