package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type CustomSelect struct {
	widget.BaseWidget
	options          []string
	selected         string
	onSelected       func(string)
	buttonStack      *fyne.Container
	buttonBackground *canvas.Rectangle
	buttonText       *canvas.Text
	button           *widget.Button
	popup            *widget.PopUp
	bgColor          map[string]color.Color
}

var _ fyne.CanvasObject = (*CustomSelect)(nil)

func NewCustomSelect(bgColor map[string]color.Color, options []string, onSelected func(string)) *CustomSelect {
	cs := &CustomSelect{
		options:    options,
		onSelected: onSelected,
		bgColor:    bgColor,
	}
	cs.ExtendBaseWidget(cs)
	// Set default BG to "almost" transparent
	bc := &color.RGBA{R: 12, G: 111, B: 211, A: 1}
	cs.buttonBackground = canvas.NewRectangle(bc)
	cs.button = widget.NewButton("", cs.showPopup)
	cs.buttonText = canvas.NewText("Select...", &colornames.White)
	cs.buttonText.Alignment = fyne.TextAlignCenter
	cs.buttonText.TextStyle = fyne.TextStyle{Bold: true}
	cs.buttonStack = container.NewStack(cs.buttonBackground, cs.buttonText, cs.button)
	// Set the button container to transparent so it shows the rectange background
	cs.button.Importance = widget.LowImportance
	cs.updateButtonText()

	return cs
}

func (cs *CustomSelect) SetBGColor(c color.Color) {
	cs.buttonBackground.FillColor = c
	cs.Refresh()
}

func (cs *CustomSelect) SetSelected(value string) {
	cs.selected = value
	cs.SetBGColor(cs.getOptionBackground(value))
	cs.updateButtonText()
	if cs.popup != nil {
		cs.popup.Hide()
	}

}

func (cs *CustomSelect) updateButtonText() {
	if cs.selected == "" {
		cs.buttonText.Text = "Select..."
	} else {
		cs.buttonText.Text = cs.selected
	}

}

func (cs *CustomSelect) showPopup() {
	if cs.popup != nil {
		cs.popup.Show()
		return
	}

	optionsContainer := container.NewVBox()
	for _, option := range cs.options {
		// Skip "emtpy" option... user needs to choose a real option
		if option == "" {
			continue
		}
		btn := cs.createOptionButton(option)
		optionsContainer.Add(btn)
	}

	scroll := container.NewVScroll(optionsContainer)
	scroll.SetMinSize(fyne.NewSize(cs.button.MinSize().Width, 239))
	content := container.NewPadded(scroll)

	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cs.button)
	buttonSize := cs.button.Size()

	popup := widget.NewPopUp(content, fyne.CurrentApp().Driver().CanvasForObject(cs.button))
	popupPos := buttonPos.Add(fyne.NewPos(0, buttonSize.Height))

	popup.ShowAtPosition(popupPos)
	cs.popup = popup
}

func (cs *CustomSelect) createOptionButton(option string) fyne.CanvasObject {
	bgColor := cs.getOptionBackground(option)
	bg := canvas.NewRectangle(bgColor)

	optionSpacer := "              "
	pt := fmt.Sprintf("%s%s%s", optionSpacer, option, optionSpacer)

	ot := canvas.NewText(pt, &colornames.White)
	ot.Alignment = fyne.TextAlignCenter
	ot.TextStyle.Bold = true

	padded := container.NewPadded(ot)
	content := container.NewStack(
		bg,
		container.NewCenter(padded),
	)

	clickable := container.NewStack(
		content,
		widget.NewButton("", func() {
			cs.SetSelected(option)
			if cs.onSelected != nil {
				cs.onSelected(option)
			}
		}),
	)

	clickable.Objects[1].(*widget.Button).Importance = widget.LowImportance

	return clickable
}

func (cs *CustomSelect) getOptionBackground(option string) color.Color {
	if cs.bgColor[option] != nil {
		return cs.bgColor[option]
	}

	return theme.Color(theme.ColorNameBackground)
}

func (cs *CustomSelect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(cs.buttonStack)
}
