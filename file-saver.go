package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func ResetSaveData() {
	SaveMatchHistory(MatchHistory{
		LP:           0,
		LPGain:       0,
		Wins:         0,
		Losses:       0,
		TotalWins:    0,
		TotalLosses:  0,
		TotalMatches: 0,
		WinRate:      0,
	})
}

func SaveMatchHistory(matchHistory MatchHistory) {
	SaveTextToFile(`results`, `wins.txt`, strconv.Itoa(matchHistory.Wins))
	SaveTextToFile(`results`, `losses.txt`, strconv.Itoa(matchHistory.Losses))
	SaveTextToFile(`results`, `lp-gain.txt`, strconv.Itoa(matchHistory.LPGain))
	SaveTextToFile(`results`, `win-rate.txt`, strconv.Itoa(matchHistory.WinRate)+`%`)
	SaveTextToFile(`results`, `lp.txt`, strconv.Itoa(matchHistory.LP))

	mhMarshalled, err := json.Marshal(&matchHistory)
	if err == nil {
		fmt.Println(string(mhMarshalled))
		SaveTextToFile(`results`, `match-history.json`, string(mhMarshalled))
	}
}

func SaveTextToFile(directory string, fileName string, text string) {
	var file *os.File
	var err error

	if directory != `` {
		err = os.Mkdir(`results`, os.FileMode(0755))
		file, err = os.Create(directory + `/` + fileName)
	} else {
		file, err = os.Create(fileName)
	}

	defer file.Close()

	_, err = file.WriteString(text)

	if err != nil {
		LogError(SaveError)
	}
}
