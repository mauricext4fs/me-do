package main

import (
	"log"
	"me-do/repository"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type UIElements struct {
	CurrentActiveTab string

	TODOTaskTable     *widget.Table
	CriticalTaskTable *widget.Table
	VeryHighTaskTable *widget.Table
	HighTaskTable     *widget.Table
	MediumTaskTable   *widget.Table
	LowTaskTable      *widget.Table
	DoneTaskTable     *widget.Table

	TODOTaskListContainer     *fyne.Container
	CriticalTaskListContainer *fyne.Container
	VeryHighTaskListContainer *fyne.Container
	HighTaskListContainer     *fyne.Container
	MediumTaskListContainer   *fyne.Container
	LowTaskListContainer      *fyne.Container
	DoneTaskListContainer     *fyne.Container

	TaskFormContainer *fyne.Container

	DBPathText            *canvas.Text
	InProgressTimerId     int64
	InProgressTimerTaskId int64
}

type TitleEditDialog struct {
	dialog.Dialog
	focusWidget fyne.Focusable
	parent      fyne.Window
}

func NewTitleEditDialog(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool), parent fyne.Window, focusOn fyne.Focusable) *TitleEditDialog {

	bD := dialog.NewForm(title, confirm, dismiss, items, callback, parent)

	d := &TitleEditDialog{
		Dialog:      bD,
		focusWidget: focusOn,
		parent:      parent,
	}
	return d
}

// ????
func (d *TitleEditDialog) Move() {

}

func (td *TODO) ShowTaskTitleEditDialog(taskRow repository.Tasks, parentWindow fyne.Window) {
	entryTitle := widget.NewEntry()
	fTitle := &widget.FormItem{
		Text:   "Title",
		Widget: entryTitle,
	}
	entryTitle.SetText(taskRow.Title)

	editTitleForm := NewTitleEditDialog(
		"Edit task: "+taskRow.Title,
		"Save",
		"Cancel",
		[]*widget.FormItem{fTitle},
		func(submited bool) {
			if submited {
				log.Println("Save was press: ", entryTitle.Text)
				log.Println("Let's save ", entryTitle.Text, " to task id: ", taskRow.ID)
				// Let's save this
				taskRow.Title = entryTitle.Text
				td.DB.UpdateTitle(taskRow.ID, entryTitle.Text)
				log.Println("Current Tab: ", td.UIElements.CurrentActiveTab)
				td.SwitchTabTo(td.UIElements.CurrentActiveTab)
				td.UIElements.CriticalTaskTable.Refresh()
			}
		},
		td.MainWindow,
		entryTitle,
	)
	editTitleForm.Show()
}

func (d *TitleEditDialog) Show() {
	d.Dialog.Show()
	if d.focusWidget != nil {
		//fyne.CurrentApp().Driver().CanvasForObject(d.FormDialog).Focus(d.focusWidget)
		go func() {
			time.Sleep(time.Millisecond * 50)
			fyne.Do(
				func() {
					d.parent.Canvas().Focus(d.focusWidget)
				})
		}()
	}
}

func (td *TODO) buildTabs() *container.AppTabs {
	todoTabContainer := td.todoTab()
	criticalTabContainer := td.criticalTab()
	veryHighTabContainer := td.veryHighTab()
	highTabContainer := td.highTab()
	mediumTabContainer := td.mediumTab()
	lowTabContainer := td.lowTab()
	doneTabContainer := td.doneTab()
	tabs := container.NewAppTabs(
		container.NewTabItem("TODO", todoTabContainer),
		container.NewTabItem("Critical", criticalTabContainer),
		container.NewTabItem("Very High", veryHighTabContainer),
		container.NewTabItem("High", highTabContainer),
		container.NewTabItem("Medium", mediumTabContainer),
		container.NewTabItem("Low", lowTabContainer),
		container.NewTabItem("Done", doneTabContainer),
		container.NewTabItem("PlaceHolder", td.getPlaceHolderFixedImage()),
	)
	td.setSwitchTabs(tabs)
	tabs.SetTabLocation(container.TabLocationTop)

	// Flag current default Tab
	td.UIElements.CurrentActiveTab = "TODO"

	// Set cmd+number shortcut for Desktop
	if _, ok := td.App.(desktop.App); ok {
		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key1,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(0)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key2,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(1)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key3,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(2)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key4,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(3)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key5,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(4)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key6,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(5)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key7,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(6)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key8,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(7)
		})

		td.MainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.Key9,
			Modifier: fyne.KeyModifierSuper,
		}, func(shortcut fyne.Shortcut) {
			tabs.SelectIndex(8)
		})

	}

	return tabs
}

func (td *TODO) SwitchTabTo(tab string) {
	switch tab {
	case "TODO":
		td.UIElements.CurrentActiveTab = "TODO"
		log.Println("TODO Tab switched!")
		td.OnTabSwitchTODO()
		// And refresh container
		td.UIElements.TODOTaskListContainer.Refresh()
	case "Critical":
		td.UIElements.CurrentActiveTab = "Critical"
		log.Println("Critical Tab switched!")
		td.OnTabSwitchCritical()
		// And refresh container
		td.UIElements.CriticalTaskListContainer.Refresh()
	case "Very High":
		td.UIElements.CurrentActiveTab = "Very High"
		log.Println("VeryHigh Tab switched!")
		td.OnTabSwitchVeryHigh()
		// And refresh container
		td.UIElements.VeryHighTaskListContainer.Refresh()
	case "High":
		td.UIElements.CurrentActiveTab = "High"
		log.Println("High Tab switched!")
		td.OnTabSwitchHigh()
		// And refresh container
		td.UIElements.HighTaskListContainer.Refresh()
	case "Medium":
		td.UIElements.CurrentActiveTab = "Medium"
		log.Println("Medium Tab switched!")
		td.OnTabSwitchMedium()
		// And refresh container
		td.UIElements.MediumTaskListContainer.Refresh()
	case "Low":
		td.UIElements.CurrentActiveTab = "Low"
		log.Println("Low Tab switched!")
		td.OnTabSwitchLow()
		// And refresh container
		td.UIElements.LowTaskListContainer.Refresh()
	case "Done":
		td.UIElements.CurrentActiveTab = "Done"
		log.Println("Done Tab switched!")
		td.OnTabSwitchDone()
		// And refresh container
		td.UIElements.DoneTaskListContainer.Refresh()
	}
}

func (td *TODO) setSwitchTabs(at *container.AppTabs) {
	at.OnSelected = func(tab *container.TabItem) {
		log.Println("Tab switching to: ", tab.Text)
		td.SwitchTabTo(tab.Text)
	}
}

func (td *TODO) buildUI() *fyne.Container {
	// Window
	td.MainWindow = td.App.NewWindow("Me Do")
	//td.MainWindow.SetFixedSize(true)
	//td.MainWindow.Resize(fyne.NewSize(750, 410))
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

	vb := container.NewVBox(
		td.UIElements.DBPathText,
		openBtn,
		saveBtn,
		td.ShowTaskForm(),
		td.getSearchContainer(),
		tabs,
	)

	openBtn.Hide()
	saveBtn.Hide()

	return vb
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

func (td *TODO) refreshPriorityTab(tabname string) {
	switch tabname {
	case "TODO":
		td.OnTabSwitchTODO()
		td.UIElements.TODOTaskListContainer.Refresh()
	case "Critical":
		td.OnTabSwitchCritical()
		td.UIElements.CriticalTaskTable.Refresh()
		td.UIElements.CriticalTaskListContainer.Refresh()
	case "Very High":
		td.OnTabSwitchVeryHigh()
		td.UIElements.VeryHighTaskListContainer.Refresh()
	case "High":
		td.OnTabSwitchHigh()
		td.UIElements.HighTaskListContainer.Refresh()
	case "Medium":
		td.OnTabSwitchMedium()
		td.UIElements.MediumTaskListContainer.Refresh()
	case "Low":
		td.OnTabSwitchLow()
		td.UIElements.LowTaskListContainer.Refresh()
	case "Done":
		td.OnTabSwitchDone()
		td.UIElements.DoneTaskListContainer.Refresh()
	}
}
