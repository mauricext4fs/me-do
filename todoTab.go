package main

import (
	"fmt"
	"log"
	"me-do/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) todoTab() *fyne.Container {
	//td.LoadTasks()
	td.Tasks = td.getTasksSlice()
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
			return len(td.Tasks), len(td.Tasks[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == 0 {
				// ID
			} else if i.Col == 3 {
				// Status
				s := widget.NewSelect(taskStatus, func(value string) {
					log.Println("Select set to ", value)
					id, _ := strconv.Atoi(td.Tasks[i.Row][0].(string))
					log.Println(id)
					td.DB.UpdateStatus(int64(id), value)
					if value == "Done" {
						//td.UIElements.TaskListContainer.RemoveAll()
						//td.LoadTasks()
					}
				})
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					s,
				}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(td.Tasks[i.Row][i.Col].(string)),
				}
			}
		})

	colWidths := []float32{1, 70, 700, 70, 70}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}

	return t
}

func (td *TODO) getTasksSlice() [][]interface{} {
	var slice [][]interface{}

	tasks, err := td.currentTODOTasks()
	if err != nil {
		log.Println("err: ", err)
	}

	slice = append(slice, []interface{}{"ID", "Position", "Title", "Priority", "Status"})

	for _, x := range tasks {
		var currentRow []interface{}

		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, strconv.FormatInt(x.Position, 10))
		currentRow = append(currentRow, fmt.Sprintf("%s", x.Title))
		currentRow = append(currentRow, fmt.Sprintf("%s", x.Priority))
		currentRow = append(currentRow, fmt.Sprintf("%s", x.Status))

		slice = append(slice, currentRow)
	}

	return slice
}

func (td *TODO) currentTODOTasks() ([]repository.Tasks, error) {
	tasks, err := td.DB.AllTODOTasks()
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}

	return tasks, nil
}
