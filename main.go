package main

import (
	"database/sql"
	"log"
	"me-do/repository"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
)

type TODO struct {
	App        fyne.App
	DB         repository.Repository
	MainWindow fyne.Window
	UIElements UIElements
	ID         int64

	Tasks     [][]interface{}
	TaskTable *widget.Table
}

func main() {
	var td TODO
	a := app.NewWithID("ch.mauricext4fs.medo")
	td.App = a

	sqlDB, err := td.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	td.setupDB(sqlDB)

	// Window
	td.MainWindow = a.NewWindow("Me Do")
	//td.MainWindow.Resize(fyne.NewSize(710, 410))
	//td.MainWindow.SetFixedSize(true)
	//td.MainWindow.CenterOnScreen()
	td.MainWindow.SetMaster()

	td.UIElements.TaskFormContainer = container.NewVBox()
	td.UIElements.TaskFormContainer.Add(td.ShowTaskForm())

	//td.UIElements.TaskListContainer = container.NewVBox()
	//td.LoadTasks()
	//td.UIElements.TaskListAdaptiveContainer = container.NewAdaptiveGrid(1, td.UIElements.TaskListContainer)
	//c.Add(td.UIElements.TaskListContainer)

	todoTab := td.todoTab()
	tabs := container.NewAppTabs(
		container.NewTabItem("TODO", todoTab),
		container.NewTabItem("PlaceHolder", td.getPlaceHolderFixedImage()),
	)
	//tabs.Refresh()
	tabs.SetTabLocation(container.TabLocationTop)
	//c.Add(tabs)

	finalContent := container.NewVBox(td.ShowTaskForm(), tabs)

	td.MainWindow.SetContent(finalContent)
	td.MainWindow.ShowAndRun()
}

func (p *TODO) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = p.App.Storage().RootURI().Path() + "/sql.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *TODO) setupDB(sqlDB *sql.DB) {
	p.DB = repository.NewSQLiteRepository(sqlDB)

	err := p.DB.Migrate()
	if err != nil {
		log.Panic(err)
	}
}
