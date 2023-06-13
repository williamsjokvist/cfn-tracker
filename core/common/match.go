package common

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type MatchHistory struct {
	CFN               string `json:"cfn"`
	LP                int    `json:"lp"`
	LPGain            int    `json:"lpGain"`
	Wins              int    `json:"wins"`
	TotalWins         int    `json:"totalWins"`
	TotalLosses       int    `json:"totalLosses"`
	TotalMatches      int    `json:"totalMatches"`
	Losses            int    `json:"losses"`
	WinRate           int    `json:"winRate"`
	Opponent          string `json:"opponent"`
	OpponentCharacter string `json:"opponentCharacter"`
	OpponentLP        int    `json:"opponentLP"`
	OpponentLeague    string `json:"opponentLeague"`
	IsWin             bool   `json:"result"`
	TimeStamp         string `json:"timestamp"`
	Date              string `json:"date"`
	WinStreak         int    `json:"winStreak"`
}

func NewMatchHistory(cfn string) *MatchHistory {
	return &MatchHistory{
		CFN:               cfn,
		LP:                0,
		LPGain:            0,
		Wins:              0,
		Losses:            0,
		TotalWins:         0,
		TotalLosses:       0,
		TotalMatches:      0,
		WinRate:           0,
		WinStreak:         0,
		IsWin:             false,
		TimeStamp:         ``,
		Date:              ``,
		Opponent:          ``,
		OpponentCharacter: ``,
		OpponentLP:        0,
		OpponentLeague:    ``,
	}
}

func (mh *MatchHistory) Log() {
	fmt.Println(`
		[`+time.Now().Format(`15:04`)+`]	
		LP:`, mh.LP, `/ 
		Gain:`, mh.LPGain, `/ 
		Wins:`, mh.Wins, `/ 
		Losses:`, mh.Losses, `/ 
		Winrate:`, mh.WinRate, `%`,
	)
}

func (mh *MatchHistory) Save() error {
	saveTextToFile(`results`, `wins.txt`, strconv.Itoa(mh.Wins))
	saveTextToFile(`results`, `losses.txt`, strconv.Itoa(mh.Losses))
	saveTextToFile(`results`, `win-rate.txt`, strconv.Itoa(mh.WinRate)+`%`)
	saveTextToFile(`results`, `win-streak.txt`, strconv.Itoa(mh.WinStreak))
	saveTextToFile(`results`, `lp.txt`, strconv.Itoa(mh.LP))
	gain := strconv.Itoa(mh.LPGain)

	if mh.LPGain > 0 {
		gain = `+` + gain
	}

	saveTextToFile(`results`, `lp-gain.txt`, gain)

	// Do not save match result if there is no opponent
	if mh.Opponent == `` {
		return nil
	}

	mhMarshalled, err := json.Marshal(&mh)
	if err != nil {
		return fmt.Errorf(`marshal match history: %w`, err)
	}

	// Save current results
	saveTextToFile(`results`, `results.json`, string(mhMarshalled))

	// Now save current results to the entire log
	var arr []MatchHistory

	pastMatches, err := os.ReadFile(`results/` + mh.CFN + `-log.json`)
	if err != nil {
		// No past matches
		str := "[" + string(mhMarshalled) + "]"
		saveTextToFile(`results`, mh.CFN+`-log.json`, str)
		return nil
	}

	err = json.Unmarshal(pastMatches, &arr)
	if err != nil {
		return fmt.Errorf(`unmarshal past match history: %v`, err)
	}

	newArr := append([]MatchHistory{*mh}, arr...)
	newArrMarshalled, err := json.Marshal(&newArr)
	if err != nil {
		return fmt.Errorf(`marshal match history: %w`, err)
	}

	saveTextToFile(`results`, mh.CFN+`-log.json`, string(newArrMarshalled))
	return nil
}

func (mh *MatchHistory) Reset() {
	cleanMh := MatchHistory{
		CFN:          ``,
		LP:           0,
		LPGain:       0,
		Wins:         0,
		Losses:       0,
		TotalWins:    0,
		TotalLosses:  0,
		TotalMatches: 0,
		WinRate:      0,
		WinStreak:    0,
		IsWin:        false,
		TimeStamp:    ``,
		Date:         ``,
	}
	cleanMh.Save()
}

func GetLog(cfn string) ([]MatchHistory, error) {
	var matchLog []MatchHistory
	pastMatches, err := os.ReadFile(fmt.Sprintf(`results/%s-log.json`, cfn))
	if err != nil {
		return nil, fmt.Errorf(`read match history: %w`, err)
	}

	json.Unmarshal(pastMatches, &matchLog)
	return matchLog, nil
}

func DeleteLog(cfn string) error {
	err := os.Remove(fmt.Sprintf(`results/%s-log.json`, cfn))
	if err != nil {
		return fmt.Errorf(`delete match history db: %w`, err)
	}
	return nil
}

func GetLoggedCFNs() ([]string, error) {
	files, err := ioutil.ReadDir(`results`)
	if err != nil {
		return nil, fmt.Errorf(`read results directory: %w`, err)
	}

	cfns := []string{}

	for _, file := range files {
		fileName := file.Name()

		if !strings.Contains(fileName, `-log.json`) {
			continue
		}

		cfn := strings.Split(fileName, `-log.json`)[0]
		cfns = append(cfns, cfn)
	}

	return cfns, nil
}

func ExportLog(cfn string) error {
	var matchHistories []MatchHistory
	pastMatches, err := os.ReadFile(fmt.Sprintf(`results/%s-log.json`, cfn))
	if err != nil {
		return fmt.Errorf(`read match history: %w`, err)
	}

	err = json.Unmarshal(pastMatches, &matchHistories)
	if err != nil {
		return fmt.Errorf(`unmarshal match history: %w`, err)
	}

	csvFile, err := os.Create(fmt.Sprintf(`results/%s.csv`, cfn))
	if err != nil {
		return fmt.Errorf(`create csv file: %w`, err)
	}

	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)

	// Header
	var header []string
	header = append(header, `Date`)
	header = append(header, `Time`)
	header = append(header, `Opponent`)
	header = append(header, `Opponent Character`)
	header = append(header, `Opponent League`)
	header = append(header, `Result`)

	writer.Write(header)

	for _, obj := range matchHistories {
		var record []string
		record = append(record, obj.Date)
		record = append(record, obj.TimeStamp)
		record = append(record, obj.Opponent)
		record = append(record, obj.OpponentCharacter)
		record = append(record, obj.OpponentLeague)

		if obj.IsWin {
			record = append(record, `W`)
		} else if !obj.IsWin {
			record = append(record, `L`)
		}

		writer.Write(record)
		record = nil
	}

	writer.Flush()
	return nil
}

func GetLastSavedMatchHistory() (*MatchHistory, error) {
	var result MatchHistory

	pastResultsFile, err := os.ReadFile(`results/results.json`)
	if err != nil {
		return nil, fmt.Errorf(`nonexisting results file: %w`, err)
	}

	err = json.Unmarshal(pastResultsFile, &result)
	if err != nil {
		return nil, fmt.Errorf(`err unmarshal match history: %w`, err)
	}

	return &result, nil
}

func saveTextToFile(directory string, fileName string, text string) {
	var file *os.File
	var err error

	var path string
	if directory != `` {
		err = os.Mkdir(`results`, os.FileMode(0755))
		if err != nil {
			fmt.Println(`err make results dir: `, err)
		}
		path = fmt.Sprintf(`%s/%s`, directory, fileName)
	} else {
		path = fileName
	}

	file, err = os.Create(path)
	if err != nil {
		fmt.Println(`err create file: `, err)
	}

	_, err = file.WriteString(text)
	defer file.Close()
	if err != nil {
		fmt.Println(`err writing to file: `, fileName)
	}
}
