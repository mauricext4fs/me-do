package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) veryHighTab() *fyne.Container {
	td.initVeryHighTab()

	td.UIElements.VeryHighTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.VeryHighTaskTable),
	)

	return td.UIElements.VeryHighTaskListContainer
}

func (td *TODO) initVeryHighTab() {
	td.LoadVeryHighTasks()
	td.UIElements.VeryHighTaskTable = td.getVeryHighTasksTable()
}

func (td *TODO) OnTabSwitchVeryHigh() {
	td.UIElements.VeryHighTaskTable = nil
	td.LoadVeryHighTasks()
	td.UIElements.VeryHighTaskTable = td.getVeryHighTasksTable()
	td.UIElements.VeryHighTaskTable.Refresh()
}

func (td *TODO) getVeryHighTasksTable() *widget.Table {

	return td.getGenericTaskTable()

}

func (td *TODO) getVeryHighStatusField(id int64, curPos int64) *CustomSelect {

	return td.getGenericStatusField(id, curPos)
}
