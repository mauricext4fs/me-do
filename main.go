package main

import (
	"database/sql"
	"log"
	"me-do/repository"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	//a.Settings().SetTheme(&MyTheme{})

	sqlDB, err := td.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	td.setupDB(sqlDB)

	finalContent := td.buildUI()

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
