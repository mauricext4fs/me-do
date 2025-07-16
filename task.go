package main

import (
	"image/color"

	"golang.org/x/image/colornames"
)

var taskPriority = []string{"", "Critical", "Very High", "High", "Medium", "Low"}
var taskStatus = []string{"Not started", "In Progress", "Paused", "Stuck", "Done"}
var taskStatusColors = map[string]color.Color{
	"Not started": &color.Gray{},                            //Grey
	"In Progress": &colornames.Orange,                       //Orange
	"Paused":      &colornames.Blue,                         //Blue
	"Stuck":       &color.NRGBA{R: 255, G: 0, B: 0, A: 255}, //Red
	"Done":        &colornames.Green,                        //Green
}
var TODOColumns = []string{"ID", "Position", "Title", "Status", "Priority"}
var TODOColumnsSize = []float32{1, 70, 600, 210, 70}

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
