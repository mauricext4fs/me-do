package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		container.NewAdaptiveGrid(1, td.TODOTaskTable),
	)

	return td.UIElements.TODOTaskListContainer
}

func (td *TODO) initTODOTab() {
	td.LoadTODOTasks()
	td.TODOTaskTable = td.getTasksTable()
}

func (td *TODO) OnTabSwitchTODO() {
	td.TODOTaskTable = nil
	td.LoadTODOTasks()
	td.TODOTaskTable = td.getTasksTable()
	td.TODOTaskTable.Refresh()
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
		td.TODOTaskTable.Refresh()
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
			id := taskRow.ID

			colName := TODODisplayColumns[i.Col]

			switch colName {
			case "Position":
				curPos := taskRow.Position

				pc := td.getActionButtonsContainer(id, curPos, taskRow.Title)

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
			log.Println("Value in cell: ", taskRow.Title)
			entryTitle := widget.NewEntry()
			fTitle := &widget.FormItem{
				Text:   "Title",
				Widget: entryTitle,
			}
			entryTitle.SetText(taskRow.Title)

			dialog.ShowForm("Edit task: "+taskRow.Title, "Save", "Cancel",
				[]*widget.FormItem{fTitle},
				func(submited bool) {
					if submited {
						log.Println("Save was press: ", entryTitle.Text)
						log.Println("Let's save ", entryTitle.Text, " to task id: ", taskRow.ID)
						// Let's save this
						taskRow.Title = entryTitle.Text
						td.DB.UpdateTitle(taskRow.ID, entryTitle.Text)
						td.LoadTODOTasks()
						td.TODOTaskTable.Refresh()
						td.UIElements.TODOTaskListContainer.Refresh()
					}
				}, td.MainWindow)
		}
	}

	return t
}

func (td *TODO) getTODOStatusField(id int64, curPos int64) *CustomSelect {
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

			// Do as if we switch the Tab and reload everything
			td.OnTabSwitchTODO()

			td.CriticalTaskTable.Refresh()

			// Refresh UI views
			td.UIElements.TODOTaskListContainer.Refresh()
			td.UIElements.CriticalTaskListContainer.Refresh()
		}

	})

	return s
}
