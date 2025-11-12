package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) todoTab() *fyne.Container {
	td.initTODOTab()

	td.UIElements.TODOTaskListContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.UIElements.TODOTaskTable),
	)

	return td.UIElements.TODOTaskListContainer
}

func (td *TODO) initTODOTab() {
	td.LoadTODOTasks()
	td.UIElements.TODOTaskTable = td.getTasksTable()
}

func (td *TODO) OnTabSwitchTODO() {
	td.UIElements.TODOTaskTable = nil
	td.LoadTODOTasks()
	td.UIElements.TODOTaskTable = td.getTasksTable()
	td.UIElements.TODOTaskTable.Refresh()
}

func (td *TODO) getSearchContainer() *fyne.Container {

	searchText := widget.NewEntry()
	searchText.SetPlaceHolder("Search TODO tasks")
	searchBtn := widget.NewButtonWithIcon("Search TODO Tasks", theme.SearchIcon(), func() {
		log.Println("Searching for: ", searchText.Text)
		res, err := td.DB.SearchTODOTasks(searchText.Text)
		if err != nil {
			log.Println("Oups... something when wrong: ", err)
		}
		td.TODOTasks = res
		td.UIElements.TODOTaskTable.Refresh()
		td.UIElements.TODOTaskListContainer.Refresh()
	})
	searchContainer := container.NewGridWithColumns(2,
		searchText,
		searchBtn,
	)

	return searchContainer
}

func (td *TODO) getTasksTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(td.TODOTasks), len(TODODisplayColumns) // Column numbers
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(" . "))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			taskRow := td.TODOTasks[i.Row]
			endRow := len(td.TODOTasks) - 1
			id := taskRow.ID

			colName := TODODisplayColumns[i.Col]

			switch colName {
			case "Position":
				curPos := taskRow.Position

				pc := td.getActionButtonsContainer(endRow, i.Row, id, curPos, taskRow.Title)

				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					pc,
				}
			case "Title":
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(taskRow.Title),
				}
			case "Status":
				// Status
				sSel := td.getTODOStatusField(id, taskRow.Position)
				sSel.SetSelected(taskRow.Status)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					sSel,
				}
			case "Priority":
				pSel := td.getPriorityField(id)
				pSel.SetSelected(taskRow.Priority)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					pSel,
				}
			case "UpdatedAt":
				// updated_at
				//tUA, _ := time.ParseDuration(taskRow.UpdatedAt.Local().GoString())
				//dUA := time.Since(taskRow.UpdatedAt)
				uAl := td.getUpdatedAtField(taskRow.UpdatedAt)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					uAl,
				}
			default:
				// Default is empty
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(""),
				}
			}

		})

	for i := 0; i < len(TODOColumnsSize); i++ {
		t.SetColumnWidth(i, TODOColumnsSize[i])
	}

	t.OnSelected = func(id widget.TableCellID) {
		log.Println("Cell was selected: ", id)
		// Which translate to:
		taskRow := td.TODOTasks[id.Row]
		colName := TODODisplayColumns[id.Col]
		log.Println("Here is the row involved: ", taskRow, " with column name: ", colName)
		if colName == "Title" {
			td.ShowTaskTitleEditDialog(taskRow, td.MainWindow)

		}
	}

	return t
}

func (td *TODO) getTODOStatusField(id int64, curPos int64) *CustomSelect {
	return td.getGenericStatusField(id, curPos, "TODO")
}
