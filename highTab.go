package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) highTab() *fyne.Container {
	td.initHighTab()

	td.UIElements.HighTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.HighTaskTable),
	)

	return td.UIElements.HighTaskListContainer
}

func (td *TODO) initHighTab() {
	td.LoadHighTasks()
	td.UIElements.HighTaskTable = td.getHighTasksTable()
}

func (td *TODO) OnTabSwitchHigh() {
	td.UIElements.HighTaskTable = nil
	td.LoadHighTasks()
	td.UIElements.HighTaskTable = td.getHighTasksTable()
	td.UIElements.HighTaskTable.Refresh()
}

func (td *TODO) getHighTasksTable() *widget.Table {

	return td.getGenericTaskTable()

}

func (td *TODO) getHightStatusField(id int64, curPos int64) *CustomSelect {

	return td.getGenericStatusField(id, curPos)
}
