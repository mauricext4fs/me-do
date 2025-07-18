package main

import (
	"log"
	"me-do/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UIElements struct {
	TaskListAdaptiveContainer *fyne.Container
	TaskListContainer         *fyne.Container
	TaskFormContainer         *fyne.Container
	TODOTasks                 []repository.Tasks
}

func (td *TODO) LoadTasks() {
	tasks, err := td.DB.AllTODOTasks()
	td.UIElements.TODOTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) buildUI() *fyne.Container {
	// Window
	td.MainWindow = td.App.NewWindow("Me Do")
	//td.MainWindow.Resize(fyne.NewSize(710, 410))
	//td.MainWindow.SetFixedSize(true)
	//td.MainWindow.CenterOnScreen()
	td.MainWindow.SetMaster()

	td.UIElements.TaskFormContainer = container.NewVBox()
	td.UIElements.TaskFormContainer.Add(td.ShowTaskForm())

	todoTab := td.todoTab()
	tabs := container.NewAppTabs(
		container.NewTabItem("TODO", todoTab),
		container.NewTabItem("PlaceHolder", td.getPlaceHolderFixedImage()),
	)
	//tabs.Refresh()
	tabs.SetTabLocation(container.TabLocationTop)
	//c.Add(tabs)

	return container.NewVBox(td.ShowTaskForm(), tabs)

}

func (td *TODO) AddTaskRow(t repository.Tasks) fyne.CanvasObject {
	hbox := container.NewHBox()
	var tr = &TaskForm{}
	tr.Position = widget.NewSelect([]string{"Up", "Down"}, func(value string) {
		var newPos int64
		if value == "Up" {
			newPos = t.Position + 1
		} else {
			newPos = t.Position - 1
		}
		td.DB.UpdatePosition(t.ID, newPos)
		log.Println("Set position to: ", newPos, " from Position: ", t.Position)
		t.Position = newPos
		td.UIElements.TaskListContainer.RemoveAll()
		td.LoadTasks()
	})
	tr.Title = widget.NewLabelWithStyle(t.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	tr.Status = widget.NewSelect(taskStatus, func(value string) {
		td.DB.UpdateStatus(t.ID, value)
		if value == "Done" {
			td.UIElements.TaskListContainer.RemoveAll()
			td.LoadTasks()
		}
	})
	tr.Status.SetSelected((t.Status))

	tr.Priority = widget.NewSelect(taskPriority, func(value string) {
		log.Println("Select set to ", value)
	})
	tr.Priority.SetSelected(t.Priority)

	tr.LastUpdate = widget.NewLabel(t.UpdatedAt.Format("2006-01-02 15:04:25"))

	hbox.Add(tr.Position)
	hbox.Add(tr.Title)
	hbox.Add(tr.Status)
	hbox.Add(tr.Priority)
	hbox.Add(tr.LastUpdate)

	return hbox
}

func (td *TODO) getPlaceHolderFixedImage() *canvas.Image {
	img := canvas.NewImageFromFile("blueblue.png")

	img.FillMode = canvas.ImageFillOriginal

	return img
}
