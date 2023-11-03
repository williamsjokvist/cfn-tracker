package data

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

type TrackingState struct {
	CFN               string `json:"cfn"`
	UserCode          string `json:"userCode"`
	LP                int    `json:"lp"`
	LPGain            int    `json:"lpGain"`
	MR                int    `json:"mr"`
	MRGain            int    `json:"mrGain"`
	Wins              int    `json:"wins"`
	TotalWins         int    `json:"totalWins"`
	TotalLosses       int    `json:"totalLosses"`
	TotalMatches      int    `json:"totalMatches"`
	Losses            int    `json:"losses"`
	WinRate           int    `json:"winRate"`
	Character         string `json:"character"`
	Opponent          string `json:"opponent"`
	OpponentCharacter string `json:"opponentCharacter"`
	OpponentLP        int    `json:"opponentLP"`
	OpponentLeague    string `json:"opponentLeague"`
	IsWin             bool   `json:"result"`
	TimeStamp         string `json:"timestamp"`
	Date              string `json:"date"`
	WinStreak         int    `json:"winStreak"`
}

const (
	logsFolder = `results/match-logs/`
)

func findLogFileFor(cfn string) (*string, error) {
	files, err := ioutil.ReadDir(logsFolder)
	if err != nil {
		return nil, fmt.Errorf(`read results directory: %w`, err)
	}

	for _, file := range files {
		fileName := file.Name()

		parts := strings.Split(fileName, `-`)
		if len(parts) != 2 {
			continue
		}

		if strings.HasPrefix(parts[1], cfn) {
			return &fileName, nil
		}
	}

	return nil, nil
}

func getLogFileNameFor(displayName string, id string) string {
	return fmt.Sprintf(`%s-%s.json`, displayName, id)
}

func NewTrackingState(cfn string) *TrackingState {
	return &TrackingState{
		CFN:               cfn,
		UserCode:          "",
		LP:                0,
		LPGain:            0,
		MR:                0,
		MRGain:            0,
		Wins:              0,
		TotalWins:         0,
		TotalLosses:       0,
		TotalMatches:      0,
		Losses:            0,
		WinRate:           0,
		Character:         ``,
		Opponent:          ``,
		OpponentCharacter: ``,
		OpponentLP:        0,
		OpponentLeague:    ``,
		IsWin:             false,
		TimeStamp:         ``,
		Date:              ``,
		WinStreak:         0,
	}
}

func (mh *TrackingState) Log() {
	fmt.Println(`
		[`+time.Now().Format(`15:04`)+`]	
		LP:`, mh.LP, `/ 
		LP Gain:`, mh.LPGain, `/ 
		MR:`, mh.MR, `/ 
		MR Gain:`, mh.MRGain, `/ 
		Wins:`, mh.Wins, `/ 
		Losses:`, mh.Losses, `/ 
		Winrate:`, mh.WinRate, `%`,
	)
}

func (mh *TrackingState) Save() error {
	saveTextToFile(`results`, `wins.txt`, strconv.Itoa(mh.Wins))
	saveTextToFile(`results`, `losses.txt`, strconv.Itoa(mh.Losses))
	saveTextToFile(`results`, `win-rate.txt`, strconv.Itoa(mh.WinRate)+`%`)
	saveTextToFile(`results`, `win-streak.txt`, strconv.Itoa(mh.WinStreak))
	saveTextToFile(`results`, `lp.txt`, strconv.Itoa(mh.LP))
	saveTextToFile(`results`, `mr.txt`, strconv.Itoa(mh.LP))
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
	var arr []TrackingState

	displayName := mh.CFN
	var id string
	if mh.UserCode != `` {
		id = mh.UserCode
	} else {
		id = mh.CFN
	}
	fileName := getLogFileNameFor(displayName, id)
	pastMatches, err := os.ReadFile(logsFolder + fileName)
	if err != nil {
		saveTextToFile(logsFolder, fileName, fmt.Sprintf(`[%s]`, string(mhMarshalled)))
		return nil
	}

	err = json.Unmarshal(pastMatches, &arr)
	if err != nil {
		return fmt.Errorf(`unmarshal past match history: %v`, err)
	}

	newArr := append([]TrackingState{*mh}, arr...)
	newArrMarshalled, err := json.Marshal(&newArr)
	if err != nil {
		return fmt.Errorf(`marshal match history: %w`, err)
	}

	saveTextToFile(logsFolder, fileName, string(newArrMarshalled))
	return nil
}

func (mh *TrackingState) Reset() {
	cleanMh := TrackingState{
		CFN:          ``,
		LP:           0,
		LPGain:       0,
		MR:           0,
		MRGain:       0,
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

func GetLog(cfn string) ([]TrackingState, error) {
	var matchLog []TrackingState

	fileName, err := findLogFileFor(cfn)
	if err != nil {
		return nil, fmt.Errorf(`error finding log file: %w`, err)
	}

	if fileName == nil {
		return nil, fmt.Errorf(`no log file found for %s`, cfn)
	}

	pastMatches, err := os.ReadFile(logsFolder + *fileName)
	if err != nil {
		return nil, fmt.Errorf(`error reading match history: %w`, err)
	}

	err = json.Unmarshal(pastMatches, &matchLog)
	if err != nil {
		return nil, err
	}

	return matchLog, nil
}

func DeleteLog(cfn string) error {
	fileName, err := findLogFileFor(cfn)
	if err != nil {
		return err
	}

	if fileName == nil {
		return fmt.Errorf(`could not find log file for %s`, cfn)
	}

	err = os.Remove(logsFolder + *fileName)
	if err != nil {
		return fmt.Errorf(`delete match history db: %w`, err)
	}

	return nil
}

func ExportLog(cfn string) error {
	var matchHistories []TrackingState

	logFileName, err := findLogFileFor(cfn)
	if err != nil || logFileName == nil {
		return fmt.Errorf(`could not find log file for %s`, cfn)
	}

	pastMatches, err := os.ReadFile(logsFolder + *logFileName)
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

func GetSavedMatchHistory(cfn string) (*TrackingState, error) {
	mh, err := GetLog(cfn)
	if err != nil {
		return nil, err
	}
	return &mh[0], nil
}

func saveTextToFile(directory string, fileName string, text string) {
	var file *os.File
	var err error

	var path string
	if directory != `` {
		_, err := os.ReadDir(directory)
		if err != nil {
			_ = os.MkdirAll(directory, os.FileMode(0755))
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
