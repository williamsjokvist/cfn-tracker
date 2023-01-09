package main

import (
	"os"
	"strconv"
)

func SaveMatchHistory(matchHistory MatchHistory) {
	SaveTextToFile(`wins.txt`, strconv.Itoa(matchHistory.wins))
	SaveTextToFile(`losses.txt`, strconv.Itoa(matchHistory.losses))
	SaveTextToFile(`lp-gain.txt`, strconv.Itoa(matchHistory.lpGain))
	SaveTextToFile(`lp.txt`, strconv.Itoa(matchHistory.lp))

}

func SaveTextToFile(fileName string, text string) {
	err := os.Mkdir(`results`, os.FileMode(0755))
	file, err := os.Create(`results/` + fileName)

	defer file.Close()

	_, err = file.WriteString(text)

	if err != nil {
		LogError(SaveErrror.message)
	}
}
