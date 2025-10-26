package main

import (
	"fmt"
	"log"
	"time"

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

func (td *TODO) buildNotesContainer(taskId int64) *fyne.Container {

	notes, err := td.DB.GetNotes(taskId)
	if err != nil {
		// Handle error
	}

	v := container.NewVBox()
	for i := range notes {
		note := notes[i]

		lText := fmt.Sprintf("Note added on: %s", note.CreatedAt.String())
		l := widget.NewLabel(lText)
		l.TextStyle.Bold = true
		v.Add(l)
		m := widget.NewRichText(
			&widget.TextSegment{Text: note.Note, Style: widget.RichTextStyleParagraph},
		)
		m.Wrapping = fyne.TextWrapWord
		//v.Add(m)
		tg := widget.NewTextGrid()
		tg.SetText(note.Note)
		v.Add(tg)

		// Show confirmation when Text is copied
		confirmBtn := widget.NewButtonWithIcon("Text Copied to clipboard!", theme.ConfirmIcon(), nil)
		confirmBtn.Importance = widget.SuccessImportance
		confirmBtn.Hide()

		// Add copy button
		cBtn := widget.NewButtonWithIcon("Copy Text", theme.ContentCopyIcon(), nil)
		cBtn.OnTapped = func() {
			clipclip := td.App.Clipboard()
			clipclip.SetContent(note.Note)

			// Show checkmark when button is clicked
			confirmBtn.Show()
			cBtn.Disable()

			go func() {
				time.Sleep(1 * time.Second)
				fyne.Do(func() {
					confirmBtn.Hide()
					cBtn.Enable()
				})

			}()

		}

		bRow := container.NewHBox(
			cBtn,
			confirmBtn,
		)

		v.Add(bRow)

	}

	return v
}
