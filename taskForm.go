package main

import (
	"log"

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

type NewTask struct {
	Title *widget.Entry
}

func (td *TODO) ShowTaskForm() fyne.CanvasObject {
	hbox := container.NewHBox()
	var tr = &TaskForm{}
	tr.Title = widget.NewLabelWithStyle("Task title", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	tr.Status = widget.NewLabel("Status")
	tr.Priority = widget.NewLabel("Priority")
	tr.LastUpdate = widget.NewLabel("Last update")

	var nt = &NewTask{}
	nt.Title = widget.NewEntry()
	nt.Title.SetPlaceHolder("Enter Task name...")

	s := widget.NewButton("Save", func() {
		log.Println("Content was: ", nt.Title.Text)
	})

	hbox.Add(nt.Title)
	hbox.Add(tr.Status)
	hbox.Add(tr.Priority)
	hbox.Add(tr.LastUpdate)
	hbox.Add(s)

	return hbox
}

func (td *TODO) ShowTaskRow() fyne.CanvasObject {
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
