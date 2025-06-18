package main

import (
	"database/sql"
	"fmt"
	"log"
	"me-do/repository"
	"os"
	"time"

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
	Countdown  Countdown
	ID         int64
	Stop       bool
}

type Countdown struct {
	Minute int64
	Second int64
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
	td.MainWindow.Resize(fyne.Size{Width: 1200, Height: 1400})
	td.MainWindow.CenterOnScreen()
	td.MainWindow.SetMaster()

	c := container.NewStack()
	c.Add(td.ShowTaskForm())
	c.Add(td.Show(c))

	td.MainWindow.SetContent(c)
	td.MainWindow.ShowAndRun()
}

func (td *TODO) Animate(co fyne.CanvasObject, win fyne.Window) {
	tick := time.NewTicker(time.Second)
	go func() {
		for !td.Stop {
			td.Layout(nil, co.Size())
			<-tick.C
			td.CountdownDown()
			td.UIElements.CountDownText.UpdateText(fmt.Sprintf("%d : %d", td.Countdown.Minute, td.Countdown.Second))
		}
		if td.Countdown.Minute == 0 && td.Countdown.Second == 0 {
			err := td.DB.UpdateActivity(td.ID, repository.Activities{ID: td.ID, EndTimestamp: time.Now()})
			if err != nil {
				log.Fatal("Error updating activity to sqlite DB: ", err)
			}

			if td.App.Preferences().FloatWithFallback("withSound", 1) == 1 {
				//PlayNotificationSound()
			}

			if td.App.Preferences().FloatWithFallback("withNotification", 0) == 1 {
				n := fyne.NewNotification("Task finished!", "Another task completed. Congrats!")
				app.New().SendNotification(n)
			}

		}
	}()
}

func (td *TODO) Reset(win fyne.Window, newTitle string) {
	// Stop any existing counter (if any)
	td.Stop = true
	time.Sleep(1 * time.Second)
	td.Countdown.Minute = 24
	td.Countdown.Second = 59
	//td.UIElements.CountDownText.UpdateText("25 Minutes")

	td.UpdateStartStopButton("Start Task", false)
	if win != nil && newTitle != "" {
		fyne.Window.SetTitle(win, newTitle)
	}
}

func (td *TODO) UpdateTODOCount() {
	//count, err := td.DB.CountCompletedTODO()
	count := 0
	/*if err != nil {
		log.Fatal("Error getting count of completed tasks from sqlite DB: ", err)
	}*/
	td.UIElements.TODOCountLabel.SetText(fmt.Sprintf("Completed Tasks: %d", count))
}

func (td *TODO) CountdownDown() {
	td.Countdown.Second--
	if td.Countdown.Minute >= 1 && td.Countdown.Second <= 0 {
		td.Countdown.Minute--
		td.Countdown.Second = 59
	} else if td.Countdown.Minute == 0 && td.Countdown.Second <= 0 {
		td.Stop = true
	}
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
