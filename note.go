package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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

	attachmentBtn := widget.NewButtonWithIcon("Attach File", theme.ContentAddIcon(), func() {
		td.ShowNotesAttachmentOpenDialog()

		// Then refresh the note list
		notesContainer.RemoveAll()
		notesContainer.Add(td.buildNotesContainer(taskId))

	})
	v.Add(attachmentBtn)

	v.Add(notesContainer)

	scroll := container.NewScroll(v)

	w.SetContent(scroll)
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
		tg := widget.NewTextGrid()
		tg.SetText(note.Note)
		v.Add(tg)

		// Show confirmation when Text is copied
		confirmBtn := widget.NewButtonWithIcon("Text Copied to clipboard!", theme.ConfirmIcon(), nil)
		// Just so we get it green
		confirmBtn.Importance = widget.SuccessImportance
		confirmBtn.Hide()

		// Add copy button
		cBtn := widget.NewButtonWithIcon("Copy Text", theme.ContentCopyIcon(), nil)
		cBtn.OnTapped = func() {
			clipclip := td.App.Clipboard()
			clipclip.SetContent(note.Note)

			// Show info when note Text is copied
			confirmBtn.Show()
			cBtn.Disable()

			go func() {
				time.Sleep(2 * time.Second)
				fyne.Do(
					func() {
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

func (td *TODO) ShowNotesAttachmentOpenDialog() {
	mainWin := td.MainWindow
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mainWin)
			return
		}

		// Nothing was choosen
		if reader == nil {
			return
		}

		defer reader.Close()

		fileURI := reader.URI()
		filename := filepath.Base(fileURI.String())
		fileExt := filepath.Ext(fileURI.String())

		// Add to DB and use the id for storage
		td.DB.AddFileToNote(1, filename, fileExt)

		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, mainWin)
			return
		}

		td.InfoLog.Println(len(data), " Data uploaded!")

	}, mainWin)
	fileDialog.Show()
	fileDialog.Resize(fyne.Size{Width: 900, Height: 1000})
	ext := []string{".jpg", ".png", ".pdf"}
	filter := storage.NewExtensionFileFilter(ext)
	fileDialog.SetFilter(filter)
}
