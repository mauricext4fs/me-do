package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"me-do/repository"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var noteFiles []repository.File
var noteWindow fyne.Window

func (td *TODO) showNotesWindow(taskId int64, taskTitle string) {

	// Create Window
	wTitle := fmt.Sprintf("Notes for task: %s", taskTitle)
	noteWindow = td.App.NewWindow(wTitle)

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
		noteId, err := td.DB.AddNote(taskId, m.Text)
		if err != nil {
			log.Println("Error saving note: ", err)
		}

		log.Println("Note Save with ID: ", noteId)

		// Save Files to Note
		for i := range noteFiles {
			file := noteFiles[i]
			fileId := file.ID

			td.InfoLog.Println("Adding fileId: ", fileId, " to noteId: ", noteId)
			_, err := td.DB.AddFileToTaskNote(fileId, noteId)
			if err != nil {
				td.InfoLog.Println("Error linking file to Note: ", err)
			}
		}

		// If all is good reset the notes text field
		m.SetText("")
		m.Refresh()
		// And the Files array
		noteFiles = nil

		// Then refresh the note list
		notesContainer.RemoveAll()
		notesContainer.Add(td.buildNotesContainer(taskId))

	})
	saveBtn.Alignment = widget.ButtonAlign(fyne.TextAlignTrailing)
	v.Add(saveBtn)

	attachmentBtn := widget.NewButtonWithIcon("Attach File", theme.ContentAddIcon(), func() {
		// Fake note Id as the note may not be saved yet at that point
		noteId := 1
		dialog := td.GetNotesAttachmentOpenDialog(int64(noteId))
		dialog.Show()
		dialog.Refresh()

		// Then refresh the note list
		notesContainer.RemoveAll()
		notesContainer.Add(td.buildNotesContainer(taskId))

	})
	v.Add(attachmentBtn)

	v.Add(notesContainer)

	scroll := container.NewScroll(v)

	noteWindow.SetContent(scroll)
	noteWindow.Resize(fyne.Size{Width: 1000, Height: 700})
	noteWindow.Show()
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

func (td *TODO) GetNotesAttachmentOpenDialog(noteId int64) *NoteFileDialog {
	mainWin := noteWindow
	fileDialog := NewNoteFileDialog(func(reader fyne.URIReadCloser, err error) {
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

		log.Println("Adding file: ", filename, " added to note id: ", 1)

		// Add to DB and use the id for storage
		fileId, err := td.DB.AddFile(filename, fileExt)
		if err != nil {
			log.Fatalln("Adding file info to the note_files table fail!  ", err)
		}
		log.Println("New note_files ID: ", fileId)

		// Replace original name with DB ID of the file
		filename = fmt.Sprintf("%s%s", strconv.Itoa(int(fileId)), fileExt)

		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, mainWin)
			return
		}

		fyneStorage := td.App.Storage().RootURI()
		fyneFileUri, err := storage.Child(fyneStorage, filename)
		if err != nil {
			td.ErrorLog.Println("Cannot create fyne URI for saving new File")
		}

		ext := []string{".jpg", ".png", ".pdf"}
		filter := storage.NewExtensionFileFilter(ext)
		if !filter.Matches(fyneFileUri) {
			dialog.ShowError(errors.New("File not supported "), mainWin)
			return
		}

		writable, err := storage.CanWrite(fyneFileUri)
		if err != nil || !writable {
			td.ErrorLog.Println("Cannot write to fyne Storage at: ", fyneFileUri)
		}

		td.InfoLog.Println(len(data), " Data uploaded!")
		td.InfoLog.Println("In Local Path: ", fyneFileUri)

		writer, err := storage.Writer(fyneFileUri)
		if err != nil {
			td.ErrorLog.Println("Fail to create storage writer: ", err)
		}

		defer writer.Close()

		_, err = io.Copy(writer, bytes.NewReader(data))
		if err != nil {
			td.ErrorLog.Println("Could not write data to Fyne storage: ", err)
		}

		// Add to global noteFiles
		fileStruct := repository.File{ID: fileId}
		noteFiles = append(noteFiles, fileStruct)

		log.Println("noteFiles: ", noteFiles)

	}, mainWin, nil)
	//fileDialog.Show()
	fileDialog.Resize(fyne.Size{Width: 700, Height: 900})

	//fileDialog.SetFilter(filter)

	return fileDialog
}
