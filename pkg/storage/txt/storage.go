package txt

import (
	"fmt"
	"os"
	"reflect"
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

func (s *Storage) SaveMatch(match model.Match) error {
	v := reflect.ValueOf(match)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Name
		value := v.Field(i)

		parsedValue := ""
		switch value.Kind() {
		case reflect.Int, reflect.Uint16:
			parsedValue = strconv.FormatInt(value.Int(), 10)
		case reflect.String:
			parsedValue = value.String()
		default:
			return fmt.Errorf("unsupported field type: %s", value.Kind())
		}

		if err := s.saveTxtFile(fmt.Sprintf("%s.txt", field), parsedValue); err != nil {
			return fmt.Errorf(`save text file: %w`, err)
		}
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
