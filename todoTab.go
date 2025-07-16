package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) todoTab() *fyne.Container {
	td.LoadTasks()
	td.TaskTable = td.getTasksTable()

	tasksTableContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, td.TaskTable),
	)

	return tasksTableContainer
}

func (td *TODO) getTasksTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(td.UIElements.TODOTasks), len(TODOColumns) // Column numbers
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			taskRow := td.UIElements.TODOTasks[i.Row]
			id := taskRow.ID

			//log.Println("Drawing row with ID: ", id, " Row ID: ", i.Row, " Col ID: ", i.Col)

			colName := TODOColumns[i.Col]
			//log.Println("Column: ", colName, " value: ", taskRow.GetValueByName(colName))

			switch colName {
			case "Position":
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(strconv.FormatInt(taskRow.Position, 10)),
				}
			case "Title":
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(taskRow.Title),
				}
			case "Priority":
				pSel := td.getPriorityField(id)
				pSel.SetSelected(taskRow.Priority)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					pSel,
				}
			case "Status":
				// Status
				sSel := td.getStatusField(id)
				sSel.SetSelected(taskRow.Status)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					sSel,
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

	return t
}
