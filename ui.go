package main

import (
	"image/color"
	"log"
	"me-do/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UIElements struct {
	TaskListAdaptiveContainer *fyne.Container
	TaskListContainer         *fyne.Container
	TaskFormContainer         *fyne.Container
	TODOTasks                 []repository.Tasks
}

type CustomSelect struct {
	widget.BaseWidget
	options    []string
	selected   string
	onSelected func(string)
	button     *widget.Button
	popup      *widget.PopUp
	bgColor    map[string]color.Color
}

var _ fyne.CanvasObject = (*CustomSelect)(nil)

func NewCustomSelect(bgColor map[string]color.Color, options []string, onSelected func(string)) *CustomSelect {
	cs := &CustomSelect{
		options:    options,
		onSelected: onSelected,
		bgColor:    bgColor,
	}
	cs.ExtendBaseWidget(cs)
	cs.button = widget.NewButton("Select...", cs.showPopup)
	cs.updateButtonText()

	return cs
}

func (cs *CustomSelect) SetSelected(value string) {
	cs.selected = value
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
	/*content := container.NewBorder(
		nil, nil, nil, nil,
		container.NewPadded(scroll),
	)*/
	content := container.NewPadded(scroll)

	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cs.button)
	buttonSize := cs.button.Size()

	popup := widget.NewPopUp(content, fyne.CurrentApp().Driver().CanvasForObject(cs.button))
	popupPos := buttonPos.Add(fyne.NewPos(0, buttonSize.Height))

	//cs.popup = widget.NewModalPopUp(content, fyne.CurrentApp().Driver().CanvasForObject(cs.button))

	//cs.popup.Show()
	popup.ShowAtPosition(popupPos)
	cs.popup = popup
}

func (cs *CustomSelect) createOptionButton(option string) fyne.CanvasObject {
	bgColor := cs.getOptionBackground(option)
	bg := canvas.NewRectangle(bgColor)

	label := widget.NewLabel(option)
	label.Alignment = fyne.TextAlignLeading

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
	return widget.NewSimpleRenderer(cs.button)
}

func (td *TODO) LoadTasks() {
	tasks, err := td.DB.AllTODOTasks()
	td.UIElements.TODOTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) drawTaskRows() {
	for _, x := range td.UIElements.TODOTasks {
		td.UIElements.TaskListContainer.Add(td.AddTaskRow(x))
		td.UIElements.TaskListContainer.Add(layout.NewSpacer())
	}
}

func (td *TODO) AddTaskRow(t repository.Tasks) fyne.CanvasObject {
	hbox := container.NewHBox()
	var tr = &TaskForm{}
	tr.Position = widget.NewSelect([]string{"Up", "Down"}, func(value string) {
		var newPos int64
		if value == "Up" {
			newPos = t.Position + 1
		} else {
			newPos = t.Position - 1
		}
		td.DB.UpdatePosition(t.ID, newPos)
		log.Println("Set position to: ", newPos, " from Position: ", t.Position)
		t.Position = newPos
		td.UIElements.TaskListContainer.RemoveAll()
		td.LoadTasks()
	})
	tr.Title = widget.NewLabelWithStyle(t.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	tr.Status = widget.NewSelect(taskStatus, func(value string) {
		td.DB.UpdateStatus(t.ID, value)
		if value == "Done" {
			td.UIElements.TaskListContainer.RemoveAll()
			td.LoadTasks()
		}
	})
	tr.Status.SetSelected((t.Status))

	tr.Priority = widget.NewSelect(taskPriority, func(value string) {
		log.Println("Select set to ", value)
	})
	tr.Priority.SetSelected(t.Priority)

	tr.LastUpdate = widget.NewLabel(t.UpdatedAt.Format("2006-01-02 15:04:25"))

	hbox.Add(tr.Position)
	hbox.Add(tr.Title)
	hbox.Add(tr.Status)
	hbox.Add(tr.Priority)
	hbox.Add(tr.LastUpdate)

	return hbox
}

func (td *TODO) getPlaceHolderFixedImage() *canvas.Image {
	img := canvas.NewImageFromFile("blueblue.png")

	img.FillMode = canvas.ImageFillOriginal

	return img
}
