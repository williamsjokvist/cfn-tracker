package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func CSV(path string, records [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create csv: %w", err)
	}
	if err := csv.NewWriter(f).WriteAll(records); err != nil {
		return errors.Join(fmt.Errorf("write csv: %w", err), f.Close())
	}
	return f.Close()
}
