package main

import (
	"database/sql"
	"io"
	"log"
	"me-do/repository"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
)

type TODO struct {
	App           fyne.App
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DefaultDBPath string
	CurrentDBPath string
	DB            repository.Repository
	MainWindow    fyne.Window
	UIElements    UIElements
	ID            int64

	Tasks     [][]interface{}
	TaskTable *widget.Table
}

func main() {
	var td TODO
	a := app.NewWithID("ch.mauricext4fs.medo")
	td.App = a
	//a.Settings().SetTheme(&MyTheme{})

	// Setting default DB location value to Fyne default
	td.DefaultDBPath = td.App.Storage().RootURI().Path() + "/sql.db"

	// Adding loggers
	td.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	td.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	path := td.getDBPath()
	sqlDB, err := td.connectSQL(path)
	if err != nil {
		log.Panic(err)
	}

	td.setupDB(sqlDB)

	finalContent := td.buildUI()

	td.MainWindow.SetContent(finalContent)
	td.MainWindow.ShowAndRun()

}

func (td *TODO) getDBPath() string {

	path := ""

	// Local overwrite for testing with Env variable
	if os.Getenv("DB_PATH") != "" {
		rpath := os.Getenv("DB_PATH")
		rpath = strings.Replace(rpath, "./", "", 1)
		cpath, err := os.Getwd()
		if err != nil {
			log.Fatalln("Cannot determine current Path: ", err)
		}
		path = cpath + "/" + rpath
		log.Println("Full path: ", path)

	} else {
		path = td.App.Preferences().StringWithFallback("currentDBPath", td.DefaultDBPath)
	}

	return path
}

func (td *TODO) connectSQL(path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Only setting current DB IF open worked!
	td.CurrentDBPath = path
	td.App.Preferences().SetString("currentDBPath", path)

	log.Println("DB: ", td.App.Preferences().String("currentDBPath"), " Opened!")

	return db, nil
}

func (td *TODO) setupDB(sqlDB *sql.DB) {
	td.DB = repository.NewSQLiteRepository(sqlDB)

	err := td.DB.Migrate()
	if err != nil {
		log.Panic(err)
	}
}

func (td *TODO) addToRecentDBFilesList(name string) []string {
	var emptyRecentPath []string
	cP := td.App.Preferences().StringListWithFallback("recentPath", emptyRecentPath)

	dub := false
	for _, v := range cP {
		if v == name {
			dub = true
		}
	}

	if dub == false {
		cP = append(cP, name)
		// Replace the recentPath
		td.App.Preferences().SetStringList("recentPath", cP)
	}

	return cP
}

func (td *TODO) copyDB(src string, dst string) error {

	log.Println("Copying DB: ", src, " to new location: ", dst)

	log.Println("Opening src DB")
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	log.Println("Opening dst DB")
	dstFile, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	log.Println("Copying...")
	_, err = io.Copy(dstFile, sourceFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
