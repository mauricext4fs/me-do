package main

import (
	"image/color"
	"log"
	"me-do/repository"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UIElements struct {
	CountDownText           *CustomText
	StartStopButton         *widget.Button
	ResetButton             *widget.Button
	QuitButton              *widget.Button
	TODOCountLabel          *widget.Label
	SoundSliderLabel        *widget.Label
	SoundSlider             *widget.Slider
	NotificationSliderLabel *widget.Label
	NotificationSlider      *widget.Slider

	TaskContainer *fyne.Container
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

func (td *TODO) LoadTasks() {
	tasks, err := td.DB.AllTasks()
	if err != nil {
		log.Println(err)
	}
	for _, x := range tasks {
		td.UIElements.TaskContainer.Add(td.AddTaskRow(x))
		td.UIElements.TaskContainer.Add(layout.NewSpacer())
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
		td.UIElements.TaskContainer.RemoveAll()
		td.LoadTasks()
	})
	tr.Title = widget.NewLabelWithStyle(t.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	tr.Status = widget.NewSelect([]string{"Not started", "In Progress", "Paused", "Stuck", "Done"}, func(value string) {
		log.Println("Select set to ", value)
		log.Println(t.ID)
		td.DB.UpdateStatus(t.ID, value)
	})
	tr.Status.SetSelected((t.Status))

	tr.Priority = widget.NewSelect([]string{"", "Critical"}, func(value string) {
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

func (td *TODO) Show(stack *fyne.Container) fyne.CanvasObject {

	content := container.NewVBox()

	td.UIElements.StartStopButton = widget.NewButton("Start ", func() {
		if td.Stop {
			result, err := td.DB.InsertTask(repository.Tasks{Title: "Task Title", CreatedAt: time.Now()})
			if err != nil {
				log.Fatal("Error adding activity to sqlite DB: ", err)
			}
			td.ID = result.ID
			fyne.Window.SetTitle(td.MainWindow, "TODO running")
			td.UpdateStartStopButton("", true)
			td.Stop = false
			go td.Animate(content, td.MainWindow)
		} else {
			fyne.Window.SetTitle(td.MainWindow, "TODO paused")
			td.UpdateStartStopButton("Continue", false)
			td.Stop = true
		}
	})
	td.UIElements.ResetButton = widget.NewButton("Reset ", func() {
		td.Reset(td.MainWindow, "MeDo")
	})
	td.UIElements.QuitButton = widget.NewButton("Quit ", func() {
		td.App.Quit()
	})

	td.UIElements.TODOCountLabel = widget.NewLabelWithStyle("Completed TODO: 0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	td.UpdateTODOCount()

	td.UIElements.SoundSliderLabel = widget.NewLabelWithStyle("Sound:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	td.UIElements.SoundSlider = widget.NewSlider(0, 1)
	td.UIElements.SoundSlider.Bind(binding.BindPreferenceFloat("withSound", td.App.Preferences()))
	td.UIElements.NotificationSliderLabel = widget.NewLabelWithStyle("Notification:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	td.UIElements.NotificationSlider = widget.NewSlider(0, 1)
	td.UIElements.NotificationSlider.Bind(binding.BindPreferenceFloat("withNotification", td.App.Preferences()))

	content.Add(layout.NewSpacer())

	content.Add(td.UIElements.StartStopButton)
	content.Add(td.UIElements.ResetButton)
	content.Add(td.UIElements.QuitButton)

	content.Add(layout.NewSpacer())
	content.Add(td.UIElements.TODOCountLabel)
	content.Add(container.New(
		layout.NewGridLayout(2),
		td.UIElements.SoundSliderLabel,
		td.UIElements.NotificationSliderLabel,
		td.UIElements.SoundSlider,
		td.UIElements.NotificationSlider))

	//td.ShowMenu()
	td.Reset(td.MainWindow, "MeDo")

	return content

}

func (td *TODO) UpdateStartStopButton(msg string, withPauseIcon bool) {
	if withPauseIcon {
		td.UIElements.StartStopButton.SetIcon(theme.MediaPauseIcon())
	} else {
		td.UIElements.StartStopButton.SetIcon(nil)

	}
	if msg == "Continue" {
		td.UIElements.StartStopButton.SetIcon(theme.MediaPlayIcon())
	} else {
		td.UIElements.StartStopButton.SetText(msg)
	}
}

func (td *TODO) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	diameter := fyne.Min(size.Width, size.Height)
	size = fyne.NewSize(diameter, diameter)
}
