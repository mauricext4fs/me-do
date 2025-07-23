package main

import (
	"database/sql"
	"log"
	"me-do/repository"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
)

type TODO struct {
	App         fyne.App
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	CurrentFile fyne.URI
	DB          repository.Repository
	MainWindow  fyne.Window
	UIElements  UIElements
	ID          int64

	Tasks     [][]interface{}
	TaskTable *widget.Table
}

func main() {
	var td TODO
	a := app.NewWithID("ch.mauricext4fs.medo")
	td.App = a
	//a.Settings().SetTheme(&MyTheme{})

	// Adding loggers
	td.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	td.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	sqlDB, err := td.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	td.setupDB(sqlDB)

	finalContent := td.buildUI()

	td.MainWindow.SetContent(finalContent)
	td.MainWindow.ShowAndRun()

}

func (td *TODO) getDBPath() string {
	dr, err := storage.Child(td.App.Storage().RootURI(), "Documents")
	du := ""
	if err != nil {
		td.ErrorLog.Println("Could not determined the Documents directory. Failing back to the default fyne.io dir")
		du = td.App.Storage().RootURI().Path()
	}

	du = dr.Path()

	return du
}

func (td *TODO) connectSQL() (*sql.DB, error) {
	path := ""

	pathi := td.getDBPath()
	td.InfoLog.Println("Path: ", pathi)

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = td.App.Storage().RootURI().Path() + "/sql.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (td *TODO) setupDB(sqlDB *sql.DB) {
	td.DB = repository.NewSQLiteRepository(sqlDB)

	err := td.DB.Migrate()
	if err != nil {
		log.Panic(err)
	}
}
