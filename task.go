package main

import (
	"log"

	"fyne.io/fyne/v2/widget"
)

var taskPriority = []string{"", "Critical", "Very High", "High", "Medium", "Low"}
var taskStatus = []string{"Not started", "In Progress", "Paused", "Stuck", "Done"}
var TODOColumns = []string{"ID", "Position", "Title", "Status", "Priority"}
var TODOColumnsSize = []float32{1, 70, 600, 210, 70}

func (td *TODO) getStatusField(id int64) *widget.Select {
	s := widget.NewSelect(taskStatus, func(value string) {
		log.Println("Select set to ", value)
		log.Println(id)
		td.DB.UpdateStatus(id, value)
		if value == "Done" {
			td.LoadTasks()
			td.TaskTable.Refresh()
		}
	})
	return s
}
