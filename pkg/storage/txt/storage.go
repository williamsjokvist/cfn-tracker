package txt

import (
	"fmt"
	"os"
	"strconv"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type Storage struct {
	directory string
}

func NewStorage() (*Storage, error) {
	directory := `results`
	err := os.MkdirAll(directory, os.FileMode(0755))
	if err != nil {
		return nil, fmt.Errorf(`create directories: %w`, err)
	}
	return &Storage{
		directory: directory,
	}, nil
}

func (s *Storage) SaveTrackingState(ts *model.TrackingState) error {
	err := s.saveTxtFile(`wins.txt`, strconv.Itoa(ts.Wins))
	if err != nil {
		return fmt.Errorf(`save wins txt: %w`, err)
	}
	err = s.saveTxtFile(`losses.txt`, strconv.Itoa(ts.Losses))
	if err != nil {
		return fmt.Errorf(`save losses txt: %w`, err)
	}
	err = s.saveTxtFile(`win-rate.txt`, strconv.Itoa(ts.WinRate)+`%`)
	if err != nil {
		return fmt.Errorf(`save win rate txt: %w`, err)
	}
	err = s.saveTxtFile(`win-streak.txt`, strconv.Itoa(ts.WinStreak))
	if err != nil {
		return fmt.Errorf(`save win streak txt: %w`, err)
	}
	err = s.saveTxtFile(`lp.txt`, strconv.Itoa(ts.LP))
	if err != nil {
		return fmt.Errorf(`save lp txt: %w`, err)
	}
	err = s.saveTxtFile(`mr.txt`, strconv.Itoa(ts.MR))
	if err != nil {
		return fmt.Errorf(`save mr txt: %w`, err)
	}
	lpGain := strconv.Itoa(ts.LPGain)
	if ts.LPGain > 0 {
		lpGain = `+` + lpGain
	}
	mrGain := strconv.Itoa(ts.MRGain)
	if ts.MRGain > 0 {
		mrGain = `+` + mrGain
	}
	err = s.saveTxtFile(`lp-gain.txt`, lpGain)
	if err != nil {
		return fmt.Errorf(`save lp gain txt: %w`, err)
	}
	err = s.saveTxtFile(`mr-gain.txt`, mrGain)
	if err != nil {
		return fmt.Errorf(`save mr gain txt: %w`, err)
	}
	return nil
}

func (s *Storage) saveTxtFile(fileName string, text string) error {
	var file *os.File
	var err error

	file, err = os.Create(fmt.Sprintf(`%s/%s`, s.directory, fileName))
	if err != nil {
		return fmt.Errorf(`create file: %w`, err)
	}
	_, err = file.WriteString(text)
	defer file.Close()
	if err != nil {
		return fmt.Errorf(`write to file: %w`, err)
	}

	return nil
}
