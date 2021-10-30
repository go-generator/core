package display

import (
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/kbinani/screenshot"
)

var errNoActiveDisplay = errors.New("no active display")

func GetActiveDisplaySize(index int) (int, int, error) {
	n := screenshot.NumActiveDisplays()
	if n < 1 {
		return 0, 0, errNoActiveDisplay
	}
	screen := screenshot.GetDisplayBounds(index)
	return screen.Dx(), screen.Dy(), nil
}

func showPopUpWindows(app fyne.App, title, message string, displaySize fyne.Size, icon fyne.Resource) {
	size := ResizeWindows(40, 30, displaySize)
	wa := app.NewWindow(title)
	if icon != nil {
		wa.SetIcon(icon)
	}
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.SetText(message)
	vScroll := container.NewScroll(entry)
	wa.SetContent(vScroll)
	if !displaySize.IsZero() {
		wa.Resize(size)
	}
	wa.CenterOnScreen()
	wa.Show()
}

func ResizeWindows(wRatio, hRatio int, size fyne.Size) fyne.Size {
	return fyne.Size{
		Width:  wRatio * size.Width / 100,
		Height: hRatio * size.Height / 100,
	}
}

func ShowErrorWindows(app fyne.App, err error, displaySize fyne.Size) {
	//red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	showPopUpWindows(app, "Error", err.Error(), fyne.Size{}, theme.HomeIcon())
}

func ShowSuccessWindows(app fyne.App, message string, displaySize fyne.Size) {
	//green := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	showPopUpWindows(app, "Success", message, fyne.Size{}, theme.HomeIcon())
}

func ShowCodeWindows(app fyne.App, code string, displaySize fyne.Size) {
	size := ResizeWindows(40, 30, displaySize)
	wa := app.NewWindow("Generated Code")
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.SetText(code)
	vScroll := container.NewScroll(entry)
	wa.SetContent(vScroll)
	if !displaySize.IsZero() {
		wa.Resize(size)
	}
	wa.CenterOnScreen()
	wa.Show()
}
