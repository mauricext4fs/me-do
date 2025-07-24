package main

import (
	"log"
	"me-do/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
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
	tabs.SetTabLocation(container.TabLocationTop)

	pt := canvas.NewText("Path: "+td.getDBPath(), colornames.Hotpink)

	openBtn := widget.NewButtonWithIcon("Open existing DB", theme.DocumentIcon(), func() {
		td.showFileOpenDialog()
		log.Println("open was clicked!")
	})

	saveBtn := widget.NewButtonWithIcon("Save as (copy current DB somewhere else)", theme.DocumentSaveIcon(), func() {
		td.showFileSaveDialog()
		log.Println("save was clicked!")
	})

	return container.NewVBox(td.ShowTaskForm(), pt, openBtn, saveBtn, tabs)

}

// TODELETE
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

func (td *TODO) showFileOpenDialog() {
	win := td.MainWindow
	saveDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
		}

		// Nothing was choosen
		if read == nil {
			return
		}

		// Add current path to recent
		td.addToRecentDBFilesList(td.CurrentDBPath)

		// Copy current DB to new Location

		// save Path
		//td.CurrentDBPath = write.URI().Path()
		log.Println("New DB Path", read.URI().Path())

		// Reset DB with new location
		db, err := td.connectSQL(read.URI().Path())
		if err != nil {
			// Not working... cannot open new DB
			log.Panicln("Cannot open new DB location ", err)
		}
		td.setupDB(db)
		td.LoadTasks()
		td.TaskTable.Refresh()

		// Add filename to the Window title
		win.SetTitle("MeDo - " + read.URI().Name())

	}, win)
	saveDialog.Show()
}

func (td *TODO) showFileSaveDialog() {
	win := td.MainWindow
	saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
		}

		// Nothing was choosen
		if write == nil {
			return
		}

		// Add current path to recent
		td.addToRecentDBFilesList(td.CurrentDBPath)

		// Copy current DB to new Location

		// save Path
		//td.CurrentDBPath = write.URI().Path()
		log.Println("New DB Path", write.URI().Path())

		// Reset DB with new location
		err = td.copyDB(td.CurrentDBPath, write.URI().Path())
		if err != nil {
			log.Panicln("Cannot copy the current DB to the new Location: ", err)
		}
		db, err := td.connectSQL(write.URI().Path())
		if err != nil {
			// Not working... cannot open new DB
			log.Panicln("Cannot open new DB location ", err)
		}
		td.setupDB(db)
		td.LoadTasks()
		td.TaskTable.Refresh()

		// Add filename to the Window title
		win.SetTitle("MeDo - " + write.URI().Name())

	}, win)
	saveDialog.Show()
}
