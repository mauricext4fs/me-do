package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (td *TODO) showNotesWindow(taskId int64, taskTitle string) {

	// Create Window
	wTitle := fmt.Sprintf("Notes for task: %s", taskTitle)
	w := td.App.NewWindow(wTitle)

	v := container.NewVBox()
	l := widget.NewLabel("New Note: ")
	v.Add(l)
	m := widget.NewMultiLineEntry()
	m.SetPlaceHolder("Write new note here")
	m.SetMinRowsVisible(10)
	v.Add(m)

	// Notes container
	notesContainer := container.NewHBox()
	notesContainer.Add(td.buildNotesContainer(taskId))

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		log.Println("Save button pressed")
		log.Println(m.Text)
		err := td.DB.AddNote(taskId, m.Text)
		if err != nil {
			log.Println("Error saving note: ", err)
		}

		// If all is good reset the notes text field
		m.SetText("")
		m.Refresh()

		// Then refresh the note list
		notesContainer.RemoveAll()
		notesContainer.Add(td.buildNotesContainer(taskId))

	})
	saveBtn.Alignment = widget.ButtonAlign(fyne.TextAlignTrailing)
	v.Add(saveBtn)
	v.Add(notesContainer)

	w.SetContent(v)
	w.Resize(fyne.Size{Width: 1000, Height: 700})
	w.Show()
}
