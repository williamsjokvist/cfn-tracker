package cmd

import (
	"errors"

	"github.com/williamsjokvist/cfn-tracker/pkg/export"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// SaveFileDialogFn returns "" if the user cancels.
type SaveFileDialogFn func(defaultFilename string) (string, error)

// Handlers share one Exporter, as the dialog needs the Wails runtime, which only
// exists once the app has started.
type Exporter struct {
	saveFileDialog SaveFileDialogFn
}

func NewExporter() *Exporter { return &Exporter{} }

func (e *Exporter) SetSaveFileDialog(saveFileDialog SaveFileDialogFn) {
	e.saveFileDialog = saveFileDialog
}

// Cancelling the dialog is not an error.
func (e *Exporter) CSV(defaultFilename string, records [][]string) error {
	if e == nil || e.saveFileDialog == nil {
		return model.WrapError(model.ErrExport, errors.New("no save file dialog"))
	}
	path, err := e.saveFileDialog(defaultFilename)
	if err != nil {
		return model.WrapError(model.ErrExport, err)
	}
	if path == "" {
		return nil
	}
	if err := export.CSV(path, records); err != nil {
		return model.WrapError(model.ErrExport, err)
	}
	return nil
}
