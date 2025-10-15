package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type UIElements struct {
	TODOTaskTable     *widget.Table
	CriticalTaskTable *widget.Table
	VeryHighTaskTable *widget.Table
	HighTaskTable     *widget.Table
	MediumTaskTable   *widget.Table
	LowTaskTable      *widget.Table

	TODOTaskListContainer     *fyne.Container
	CriticalTaskListContainer *fyne.Container
	VeryHighTaskListContainer *fyne.Container
	HighTaskListContainer     *fyne.Container
	MediumTaskListContainer   *fyne.Container
	LowTaskListContainer      *fyne.Container
	TaskFormContainer         *fyne.Container
	DBPathText                *canvas.Text
	InProgressTimerId         int64
	InProgressTimerTaskId     int64
}

func (td *TODO) buildTabs() *container.AppTabs {
	todoTabContainer := td.todoTab()
	criticalTabContainer := td.criticalTab()
	tabs := container.NewAppTabs(
		container.NewTabItem("TODO", todoTabContainer),
		container.NewTabItem("Status: Critical", criticalTabContainer),
		container.NewTabItem("PlaceHolder", td.getPlaceHolderFixedImage()),
	)
	td.setSwitchTabs(tabs)
	tabs.SetTabLocation(container.TabLocationTop)

	return tabs
}

func (td *TODO) setSwitchTabs(at *container.AppTabs) {
	at.OnSelected = func(tab *container.TabItem) {
		log.Println("Tab switching to: ", tab.Text)
		switch tab.Text {
		case "Status: Critical":
			// Load Critical tabs data
			log.Println("Critical Tab switched!")
			td.OnTabSwitchCritical()
			// And refresh container
			td.UIElements.CriticalTaskListContainer.Refresh()
		case "TODO":
			log.Println("TODO Tab switched!")
			td.OnTabSwitchTODO()
			// And refresh container
			td.UIElements.TODOTaskListContainer.Refresh()
		}
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

	tabs := td.buildTabs()

	pt := canvas.NewText("Path: "+td.CurrentDBPath, colornames.Hotpink)
	td.UIElements.DBPathText = pt

	openBtn := widget.NewButtonWithIcon("Open existing DB", theme.DocumentIcon(), func() {
		td.showFileOpenDialog()
		log.Println("open was clicked!")
	})

	saveBtn := widget.NewButtonWithIcon("Save as (copy current DB somewhere else)", theme.DocumentSaveIcon(), func() {
		td.showFileSaveDialog()
		log.Println("save was clicked!")
	})

	return container.NewVBox(td.ShowTaskForm(), td.UIElements.DBPathText, openBtn, saveBtn, td.getSearchContainer(), tabs)

}

func (td *TODO) getPlaceHolderFixedImage() *canvas.Image {
	img := canvas.NewImageFromResource(resourceBluebluePng)

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
		td.LoadTODOTasks()
		td.UIElements.TODOTaskTable.Refresh()

		// Add filename to the Window title
		win.SetTitle("MeDo - " + read.URI().Name())
		td.UIElements.DBPathText.Text = read.URI().Path()

	}, win)
	saveDialog.Show()
	saveDialog.Resize(fyne.Size{Width: 900, Height: 1000})
	ext := []string{".medo"}
	filter := storage.NewExtensionFileFilter(ext)
	saveDialog.SetFilter(filter)
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
		td.LoadTODOTasks()
		td.UIElements.TODOTaskTable.Refresh()

		// Add filename to the Window title
		win.SetTitle("MeDo - " + write.URI().Name())

	}, win)
	saveDialog.Show()
}
