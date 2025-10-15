package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) criticalTab() *fyne.Container {
	td.initCriticalTab()

	td.UIElements.CriticalTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.CriticalTaskTable),
	)

	return td.UIElements.CriticalTaskListContainer
}

func (td *TODO) initCriticalTab() {
	td.LoadCriticalTasks()
	td.UIElements.CriticalTaskTable = td.getCriticalTasksTable()
}

func (td *TODO) OnTabSwitchCritical() {
	td.UIElements.CriticalTaskTable = nil
	td.LoadCriticalTasks()
	td.UIElements.CriticalTaskTable = td.getCriticalTasksTable()
	td.UIElements.CriticalTaskTable.Refresh()
}

func (td *TODO) getCriticalTasksTable() *widget.Table {

	return td.getGenericTaskTable()

}

func (td *TODO) getCriticalStatusField(id int64, curPos int64) *CustomSelect {
	s := NewCustomSelect(taskStatusColors, taskStatus, func(value string) {
		// Stop any existing timer
		previousStatus, _ := td.DB.GetStatusByTaskID(id)

		log.Println("Previous Status: ", previousStatus, "New Status: ", value, " for taskId: ", id)

		if previousStatus == "In Progress" {
			log.Println("Previous Status was 'In Progress'; We need to stop the previous timer; if any.")
			activeTimerID, err := td.DB.GetActiveTimerByTaskId(id)
			if err != nil {
				log.Println(err)
			}
			log.Println("Active Timer ID: ", activeTimerID)
			if activeTimerID != 0 {
				td.DB.StopTimer(activeTimerID)
			}
		}

		if value == "In Progress" {
			log.Println("In Progress...")
			// Need to stop any existing timer and switch their status automatically
			if td.UIElements.InProgressTimerId != 0 {
				// ...
			}

			timer, err := td.DB.StartTimer(id)
			if err != nil {
				log.Println(err)
			}

			td.UIElements.InProgressTimerTaskId = timer.TaskID
			td.UIElements.InProgressTimerId = timer.ID

		}

		td.DB.UpdateStatus(id, value)

		if value == "Done" {
			// We need to unshift the position
			log.Println("Shifting task id: ", id, " with position ", curPos)
			td.DB.ShiftPosition(id, curPos, "TODO")
			td.LoadCriticalTasks()
			// Do we still need that CriticalTaskTable.Refresh()??
			td.UIElements.CriticalTaskTable.Refresh()

			// Refresh UI views
			td.UIElements.TODOTaskListContainer.Refresh()
			td.UIElements.CriticalTaskListContainer.Refresh()
		}

	})

	return s
}
