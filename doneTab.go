package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) doneTab() *fyne.Container {
	td.initDoneTab()

	td.UIElements.DoneTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.DoneTaskTable),
	)

	return td.UIElements.DoneTaskListContainer
}

func (td *TODO) initDoneTab() {
	td.LoadDoneTasks()
	td.UIElements.DoneTaskTable = td.getDoneTasksTable()
}

func (td *TODO) OnTabSwitchDone() {
	td.UIElements.DoneTaskTable = nil
	td.LoadDoneTasks()
	td.UIElements.DoneTaskTable = td.getDoneTasksTable()
	td.UIElements.DoneTaskTable.Refresh()
}

func (td *TODO) getDoneTasksTable() *widget.Table {

	return td.getGenericTaskTable("Done")

}

func (td *TODO) getDoneStatusField(id int64, curPos int64) *CustomSelect {

	return td.getGenericStatusField(id, curPos, "Done")
}
