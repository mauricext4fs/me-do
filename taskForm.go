package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TaskForm struct {
	Title      *widget.Label
	Status     *widget.Label
	Priority   *widget.Label
	LastUpdate *widget.Label
}

func (td *TODO) ShowTaskForm() fyne.CanvasObject {
	hbox := container.NewHBox()
	var tr = &TaskForm{}
	tr.Title = widget.NewLabelWithStyle("Task title", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	tr.Status = widget.NewLabel("Status")
	tr.Priority = widget.NewLabel("Priority")
	tr.LastUpdate = widget.NewLabel("Last update")
	hbox.Add(tr.Title)
	hbox.Add(tr.Status)
	hbox.Add(tr.Priority)
	hbox.Add(tr.LastUpdate)

	return hbox
}
