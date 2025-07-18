package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomSelect struct {
	widget.BaseWidget
	options          []string
	selected         string
	onSelected       func(string)
	buttonStack      *fyne.Container
	buttonBackground *canvas.Rectangle
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
	cs.button = widget.NewButton("Select...", cs.showPopup)
	cs.buttonStack = container.NewStack(cs.buttonBackground, cs.button)
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
		cs.button.SetText("Select....")
	} else {
		cs.button.SetText(cs.selected)
	}

}

func (cs *CustomSelect) showPopup() {
	if cs.popup != nil {
		cs.popup.Show()
		return
	}

	optionsContainer := container.NewVBox()
	for _, option := range cs.options {
		btn := cs.createOptionButton(option)
		optionsContainer.Add(btn)
	}

	scroll := container.NewVScroll(optionsContainer)
	scroll.SetMinSize(fyne.NewSize(cs.button.MinSize().Width, 300))
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

	label := widget.NewLabel(option)
	label.Alignment = fyne.TextAlignLeading
	//label.TextStyle.Bold = true

	paddedContent := container.NewHBox(layout.NewSpacer(), label, layout.NewSpacer())
	content := container.NewPadded(paddedContent)

	stack := container.NewStack(bg, content)

	clickable := container.NewStack(
		stack,
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
