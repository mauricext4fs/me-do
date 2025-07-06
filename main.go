package main

import (
	"database/sql"
	"log"
	"me-do/repository"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	_ "github.com/glebarez/go-sqlite"
)

type TODO struct {
	App        fyne.App
	DB         repository.Repository
	MainWindow fyne.Window
	UIElements UIElements
	ID         int64
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
	td.MainWindow.Resize(fyne.NewSize(600, 410))
	//	td.MainWindow.SetFixedSize(true)
	//td.MainWindow.CenterOnScreen()
	td.MainWindow.SetMaster()

	c := container.NewVBox()
	//c.Resize(fyne.Size{Width: 1000, Height: 20})

	td.UIElements.TaskFormContainer = container.NewVBox()
	td.UIElements.TaskFormContainer.Add(td.ShowTaskForm())
	c.Add(td.UIElements.TaskFormContainer)

	td.UIElements.TaskListContainer = container.NewVBox()
	td.LoadTasks()
	td.UIElements.TaskListAdaptiveContainer = container.NewAdaptiveGrid(1, td.UIElements.TaskListContainer)

	c.Add(td.UIElements.TaskListContainer)

	td.MainWindow.SetContent(c)
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
