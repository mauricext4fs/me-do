package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) mediumTab() *fyne.Container {
	td.initMediumTab()

	td.UIElements.MediumTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.MediumTaskTable),
	)

	return td.UIElements.MediumTaskListContainer
}

func (td *TODO) initMediumTab() {
	td.LoadMediumTasks()
	td.UIElements.MediumTaskTable = td.getMediumTasksTable()
}

func (td *TODO) OnTabSwitchMedium() {
	td.UIElements.MediumTaskTable = nil
	td.LoadMediumTasks()
	td.UIElements.MediumTaskTable = td.getMediumTasksTable()
	td.UIElements.MediumTaskTable.Refresh()
}

func (td *TODO) getMediumTasksTable() *widget.Table {

	return td.getGenericTaskTable()

}

func (td *TODO) getMediumStatusField(id int64, curPos int64) *CustomSelect {

	return td.getGenericStatusField(id, curPos)
}
