package main

import (
	"log"
	"me-do/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TaskForm struct {
	Title      *widget.Label
	Status     *widget.Label
	Priority   *widget.Select
	LastUpdate *widget.Label
}

type NewTask struct {
	Title *widget.Entry
}

func (td *TODO) ShowTaskForm() fyne.CanvasObject {
	hbox := container.NewHBox()
	var tr = &TaskForm{}
	tr.Status = widget.NewLabel("Status")
	tr.Priority = widget.NewSelect([]string{"", "Critical"}, func(value string) {
		log.Println("Select set to ", value)
	})

	tr.LastUpdate = widget.NewLabel("Last update")

	var nt = &NewTask{}
	nt.Title = widget.NewEntry()
	nt.Title.SetPlaceHolder("Enter Task name...")

	s := widget.NewButton("Save", func() {
		log.Println("Content was: ", nt.Title.Text)
		_, err := td.DB.InsertTask(repository.Tasks{
			Title:    nt.Title.Text,
			Priority: tr.Priority.Selected,
		})
		if err != nil {
			log.Println(err)
		}

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
	tr.Priority = widget.NewSelect([]string{"", "Critical"}, func(value string) {
		log.Println("Select set to ", value)
	})

	tr.LastUpdate = widget.NewLabel("Last update")

	hbox.Add(tr.Title)
	hbox.Add(tr.Status)
	hbox.Add(tr.Priority)
	hbox.Add(tr.LastUpdate)

	return hbox
}
