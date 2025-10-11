package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) criticalTab() *fyne.Container {
	//td.TaskTable = td.getCriticalTasksTable()
	td.OnTabSwitchCritical()

	tasksTableContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.TaskTable),
	)

	return tasksTableContainer
}

func (td *TODO) OnTabSwitchCritical() {
	td.TaskTable = nil
	td.LoadCriticalTasks()
	td.TaskTable = td.getCriticalTasksTable()
	td.TaskTable.Refresh()
}

func (td *TODO) getCriticalTasksTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(td.CriticalTasks), len(TODODisplayColumns) // Column numbers
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(" . "))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			taskRow := td.CriticalTasks[i.Row]
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
				sSel := td.getCriticalStatusField(id, taskRow.Position)
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
		taskRow := td.CriticalTasks[id.Row]
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
						td.LoadCriticalTasks()
						td.TaskTable.Refresh()
					}
				}, td.MainWindow)
		}
	}

	return t
}

func (td *TODO) getCriticalStatusField(id int64, curPos int64) *CustomSelect {
	s := NewCustomSelect(taskStatusColors, taskStatus, func(value string) {
		log.Println("Status: ", value, " for taskId: ", id)
		// Stop any existing timer
		previousStatus, _ := td.DB.GetStatusByTaskID(id)
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
			td.TaskTable.Refresh()
		}
	})

	return s
}
