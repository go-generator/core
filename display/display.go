package display

import (
	"errors"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

var errNoActiveDisplay = errors.New("no active display")

func Notify(content string, title string) {
	fyne.CurrentApp().SendNotification(fyne.NewNotification(title, content))
}

func SetIcon(path string) (fyne.Resource, error) {
	settingIcon, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	r2, err1 := fyne.LoadResourceFromPath(settingIcon)
	if err1 != nil {
		return nil, err1
	}
	return r2, nil
}

func GetActiveDisplaySize(index int) (int, int, error) {
	n := screenshot.NumActiveDisplays()
	if n < 1 {
		return 0, 0, errNoActiveDisplay
	}
	screen := screenshot.GetDisplayBounds(index)
	return screen.Dx(), screen.Dy(), nil
}

func PopUpWindows(message string, c fyne.Canvas) {
	var popup *widget.PopUp
	popup = widget.NewModalPopUp(container.NewBorder(nil, widget.NewButton("Close", func() {
		if popup != nil {
			popup.Hide()
		}
	}), nil, nil, widget.NewLabel(message)), c)
	popup.Show()
}

func ResizeWindows(wRatio, hRatio float32, size fyne.Size) fyne.Size {
	return fyne.Size{
		Width:  wRatio * size.Width / 100,
		Height: hRatio * size.Height / 100,
	}
}
