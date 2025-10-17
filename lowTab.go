package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) lowTab() *fyne.Container {
	td.initLowTab()

	td.UIElements.LowTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.LowTaskTable),
	)

	return td.UIElements.LowTaskListContainer
}

func (td *TODO) initLowTab() {
	td.LoadLowTasks()
	td.UIElements.LowTaskTable = td.getLowTasksTable()
}

func (td *TODO) OnTabSwitchLow() {
	td.UIElements.LowTaskTable = nil
	td.LoadLowTasks()
	td.UIElements.LowTaskTable = td.getLowTasksTable()
	td.UIElements.LowTaskTable.Refresh()
}

func (td *TODO) getLowTasksTable() *widget.Table {

	return td.getGenericTaskTable("Low")

}

func (td *TODO) getLowStatusField(id int64, curPos int64) *CustomSelect {

	return td.getGenericStatusField(id, curPos, "Low")
}
