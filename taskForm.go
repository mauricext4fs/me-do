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
		log.Println("Task Priority for new task set to ", value)
	})

	tr.LastUpdate = widget.NewLabel("Last update")

	var nt = &NewTask{}
	nt.Title = widget.NewEntry()
	nt.Title.SetPlaceHolder("Enter Task name...")

	s := widget.NewButton("Add new Task", func() {
		log.Println("Content was: ", nt.Title.Text)
		task := repository.Tasks{
			Title:    nt.Title.Text,
			Priority: tr.Priority.Selected,
		}
		nTask, err := td.DB.InsertTask(task)
		if err != nil {
			log.Println(err)
		}

		log.Println("New Task Inserted: ", nTask)

		// Reload the Tabs Table
		td.OnTabSwitchCritical()
		td.OnTabSwitchTODO()
		td.UIElements.TODOTaskListContainer.Refresh()
		td.UIElements.CriticalTaskListContainer.Refresh()

		//Clear up existing field value
		nt.Title.Text = ""
		nt.Title.Refresh()
		td.MainWindow.Canvas().Focus(nt.Title)

		tr.Priority.ClearSelected()
	})

	taskForm := container.NewGridWithColumns(3,
		nt.Title,
		tr.Priority,
		s,
	)

	return taskForm
}
