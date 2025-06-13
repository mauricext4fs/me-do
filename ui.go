package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type UIElements struct {
	StartStopButton         *widget.Button
	ResetButton             *widget.Button
	QuitButton              *widget.Button
	SoundSliderLabel        *widget.Label
	SoundSlider             *widget.Slider
	NotificationSliderLabel *widget.Label
	NotificationSlider      *widget.Slider
}

type CustomText struct {
	canvas.Text
}

var _ fyne.CanvasObject = (*CustomText)(nil)

func NewCustomText(text string, c color.Color) *CustomText {
	size := fyne.CurrentApp().Settings().Theme().Size("custom_text")
	nct := &CustomText{}
	nct.Text.Text = text
	nct.Text.TextSize = size
	nct.Text.Color = c

	return nct
}

func (t *CustomText) UpdateText(text string) {
	t.Text.Text = text
	t.Text.Refresh()
}
