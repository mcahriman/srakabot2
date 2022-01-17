package main

import (
	_ "embed"
	"fmt"

	"github.com/gocarina/gocsv"
)

//go:embed resources/testing/spamrank.csv
var spamranksCsv string

//go:embed resources/testing/karmarank.csv
var karmarankCsv string

//rankName,emoji,picture,minMessages,congratulationMessage,showRankCongrats

type SpamRank struct {
	ChatId                  string `csv:"-" json:"chatId"`
	RankName                string `csv:"rankName" json:"rankName"`
	Emoji                   string `csv:"emoji" json:"emoji"`
	Picture                 string `csv:"picture" json:"picture"`
	MinMessages             int    `csv:"minMessages" json:"minMessages"`
	CongratulationMessage   string `csv:"congratulationMessage" json:"congratulationMessage"`
	ShowRankCongratsTimeout int    `csv:"showRankCongrats" json:"showRankCongrats"`
}

type KarmaRank struct {
	Karma int    `csv:"karma" json:"karma"`
	Rank  string `csv:"rank" json:"rank"`
}

func getRanksCSVEmbeded() []*SpamRank {
	spamRanks := []*SpamRank{}
	err := gocsv.UnmarshalString(spamranksCsv, &spamRanks)
	fmt.Printf("%+v, err: %+v", spamRanks, err)
	return spamRanks
}

func getKarmaRanksCsvEmbedded() []*KarmaRank {

	karmaRanks := []*KarmaRank{}
	err := gocsv.UnmarshalString(karmarankCsv, &karmaRanks)
	if err != nil {
		fmt.Printf("%+v, err: %+v", karmaRanks, err)

	}
	return karmaRanks
}

func calculateDesignation(karma int, spam int) string {
	kRank := findKRankByVal(karma)
	sRank := findRankByVal(spam)
	return fmt.Sprintf("%s %s %s", sRank.Emoji, kRank.Rank, sRank.RankName)
}

func findKRankByVal(val int) KarmaRank {
	nextRank := KarmaRank{Karma: 1000}
	for _, r := range getKarmaRanksCsvEmbedded() {
		if r.Karma > nextRank.Karma && r.Karma <= val {
			nextRank = *r
		}
	}
	return nextRank
}

func findRankByVal(val int) SpamRank {
	nextRank := SpamRank{MinMessages: -900000}
	for _, r := range getRanksCSVEmbeded() {
		if r.MinMessages > nextRank.MinMessages && r.MinMessages <= val {
			nextRank = *r
		}
	}
	return nextRank
}

func getRanksDescriptionsSplitted() [][]*SpamRank {
	ranks := getRanksCSVEmbeded()
	rankChunks := [][]*SpamRank{}
	for i := 0; i < len(ranks); i += 15 {
		if i+15 < len(ranks) {
			rankChunks = append(rankChunks, ranks[i:i+15])
		} else {
			rankChunks = append(rankChunks, ranks[i:])
		}
	}
	return rankChunks
}

func getKarmaRanksDescriptionsSplitted() [][]*KarmaRank {
	ranks := getKarmaRanksCsvEmbedded()
	rankChunks := [][]*KarmaRank{}
	for i := 0; i < len(ranks); i += 15 {
		if i+15 < len(ranks) {
			rankChunks = append(rankChunks, ranks[i:i+15])
		} else {
			rankChunks = append(rankChunks, ranks[i:])
		}
	}
	return rankChunks
}
