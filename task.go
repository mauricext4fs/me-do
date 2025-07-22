package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

var taskPriority = []string{"", "Critical", "Very High", "High", "Medium", "Low"}
var taskPriorityColors = map[string]color.Color{
	"Low":       &colornames.Lightcyan,    // Cyan
	"Medium":    &colornames.Lightskyblue, // Medium blue
	"High":      &colornames.Mediumblue,   //Blue
	"Very High": &colornames.Darkblue,     //Dark Blue
	"Critical":  &colornames.Black,        //Black
}
var taskStatus = []string{"Not started", "In Progress", "Paused", "Stuck", "Done"}
var taskStatusColors = map[string]color.Color{
	"Not started": &colornames.Lightgray,                       //Grey
	"In Progress": &colornames.Orange,                          //Orange
	"Paused":      &colornames.Dodgerblue,                      //Blue
	"Stuck":       &color.NRGBA{R: 205, G: 65, B: 79, A: 255},  //Red
	"Done":        &color.NRGBA{R: 90, G: 197, B: 125, A: 255}, //Green
}
var TODODisplayColumns = []string{"Position", "Title", "Status", "Priority"}

// var TODOColumns = []string{"ID", "Position", "Title", "Status", "Priority"}
var TODOColumnsSize = []float32{80, 600, 210, 210}

//var TODOColumnsSize = []float32{1, 70, 600, 210, 210}

func (td *TODO) getStatusField(id int64) *CustomSelect {
	s := NewCustomSelect(taskStatusColors, taskStatus, func(value string) {
		td.DB.UpdateStatus(id, value)
		if value == "Done" {
			td.LoadTasks()
			td.TaskTable.Refresh()
		}
	})

	return s
}

func (td *TODO) getPriorityField(id int64) *CustomSelect {
	s := NewCustomSelect(taskPriorityColors, taskPriority, func(value string) {
		td.DB.UpdatePriority(id, value)
	})
	return s
}

func (td *TODO) getUpDownPositionField(id int64, curPos int64) *fyne.Container {
	upBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		td.DB.UpdatePosition(id, (curPos + 1))
		td.LoadTasks()
		td.TaskTable.Refresh()
	})
	downBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		td.DB.UpdatePosition(id, (curPos - 1))
		td.LoadTasks()
		td.TaskTable.Refresh()
	})
	pc := container.NewCenter(container.NewHBox(downBtn, upBtn))

	return pc
}
