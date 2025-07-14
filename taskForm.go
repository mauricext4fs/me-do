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
	Position   *widget.Select
	Status     *widget.Select
	Priority   *widget.Select
	LastUpdate *widget.Label
}

type NewTask struct {
	Title *widget.Entry
}

func (td *TODO) ShowTaskForm() fyne.CanvasObject {
	var tr = &TaskForm{}
	tr.Priority = widget.NewSelect(taskPriority, func(value string) {
		log.Println("Select set to ", value)
	})

	tr.LastUpdate = widget.NewLabel("Last update")

	var nt = &NewTask{}
	nt.Title = widget.NewEntry()
	nt.Title.SetPlaceHolder("Enter Task name...")

	s := widget.NewButton("Add new Task", func() {
		log.Println("Content was: ", nt.Title.Text)
		_, err := td.DB.InsertTask(repository.Tasks{
			Title:    nt.Title.Text,
			Priority: tr.Priority.Selected,
		})
		if err != nil {
			log.Println(err)
		}

		// Reload the TODO Table
		td.LoadTasks()
		td.TaskTable.Refresh()

		//Clear up existing field value
		nt.Title.Text = ""
		nt.Title.Refresh()

		tr.Priority.ClearSelected()
	})

	taskForm := container.NewGridWithColumns(3,
		nt.Title,
		tr.Priority,
		s,
	)

	return taskForm
}
