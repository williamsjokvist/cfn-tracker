package file

import (
	"fmt"
	"os"
)

func SaveToFile(directory string, fileName string, str string) {
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

	_, err = file.WriteString(str)
	defer file.Close()
	if err != nil {
		fmt.Println(`err writing to file: `, fileName)
	}
}
