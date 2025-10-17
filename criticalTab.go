package main

import (
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

	return td.getGenericTaskTable("Critical")

}
