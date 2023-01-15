package main

import (
	"os"
	"strconv"
)

func SaveMatchHistory(matchHistory MatchHistory) {
	SaveTextToFile(`results`, `wins.txt`, strconv.Itoa(matchHistory.wins))
	SaveTextToFile(`results`, `losses.txt`, strconv.Itoa(matchHistory.losses))
	SaveTextToFile(`results`, `lp-gain.txt`, strconv.Itoa(matchHistory.lpGain))
	SaveTextToFile(`results`, `win-rate.txt`, strconv.Itoa(matchHistory.winrate)+`%`)
	SaveTextToFile(`results`, `lp.txt`, strconv.Itoa(matchHistory.lp))
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
