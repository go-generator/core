package list

import (
	"fyne.io/fyne/v2/data/binding"
	metadata "github.com/go-generator/core"
	"path/filepath"
	"strconv"
)

func ShowFiles(data binding.ExternalStringList, dataSt binding.Struct, files []metadata.File) error {
	err := data.Set(nil)
	if err != nil {
		return err
	}
	for i := range files {
		filename := strconv.Itoa(i+1) + ". " + filepath.Base(files[i].Name)
		err = data.Append(filename)
		if err != nil {
			return err
		}
		err = dataSt.SetValue(filename, files[i].Content)
		if err != nil {
			return err
		}
	}
	return nil
}
