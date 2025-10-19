package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

var taskPriority = []string{"", "Critical", "Very High", "High", "Medium", "Low"}
var taskPriorityColors = map[string]color.Color{
	"Low":       &colornames.Lightskyblue,    // Medium blue
	"Medium":    &colornames.Mediumslateblue, // Slightly violet
	"High":      &colornames.Mediumblue,      //Blue
	"Very High": &colornames.Darkblue,        //Dark Blue
	"Critical":  &colornames.Black,           //Black
}
var taskStatus = []string{"Not started", "In Progress", "Paused", "Stuck", "Done"}
var taskStatusColors = map[string]color.Color{
	"Not started": &colornames.Lightgray,                       //Grey
	"In Progress": &colornames.Orange,                          //Orange
	"Paused":      &colornames.Dodgerblue,                      //Blue
	"Stuck":       &color.NRGBA{R: 205, G: 65, B: 79, A: 255},  //Red
	"Done":        &color.NRGBA{R: 90, G: 197, B: 125, A: 255}, //Green
}
var TODODisplayColumns = []string{"Position", "Title", "Status", "Priority", "UpdatedAt"}
var TODOColumnsSize = []float32{110, 600, 210, 210, 180}

func (td *TODO) getPriorityField(id int64) *CustomSelect {
	s := NewCustomSelect(taskPriorityColors, taskPriority, func(value string) {
		td.DB.UpdatePriority(id, value)
	})
	return s
}

func (td *TODO) getActionButtonsContainer(id int64, curPos int64, title string) *fyne.Container {
	upBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		err := td.DB.UpPosition(id, (curPos), "TODO")
		if err != nil {
			log.Println("Error on Move up press for TaskID: ", id, " : ", err)
		}
		//td.LoadTODOTasks()
		//td.TODOTaskTable.Refresh()
		// Do as if we switch the Tab and reload everything
		td.OnTabSwitchCritical()
		td.OnTabSwitchTODO()

		td.UIElements.CriticalTaskTable.Refresh()
		td.UIElements.TODOTaskTable.Refresh()

		// Refresh UI views
		td.UIElements.TODOTaskListContainer.Refresh()
		td.UIElements.CriticalTaskListContainer.Refresh()
	})
	downBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		err := td.DB.DownPosition(id, (curPos), "TODO")
		if err != nil {
			log.Println("Error on Move Down press for TaskID: ", id, " : ", err)
		}
		//td.LoadTODOTasks()
		//td.TODOTaskTable.Refresh()

		// Do as if we switch the Tab and reload everything
		td.OnTabSwitchCritical()
		td.OnTabSwitchTODO()

		td.UIElements.CriticalTaskTable.Refresh()
		td.UIElements.TODOTaskTable.Refresh()

		// Refresh UI views
		td.UIElements.TODOTaskListContainer.Refresh()
		td.UIElements.CriticalTaskListContainer.Refresh()

	})
	notesBtn := widget.NewButtonWithIcon("", theme.FileTextIcon(), func() {
		// Load notes (if any) and show them somewhere
		//err := td.DB.DownPosition(id, (curPos), "TODO")

		td.showNotesWindow(id, title)
	})
	pc := container.NewCenter(container.NewHBox(downBtn, upBtn, notesBtn))

	return pc
}

func (td *TODO) getUpdatedAtField(uT time.Time) *widget.Label {

	dUA := time.Since(uT)
	tSeconds := int(dUA.Abs().Seconds())

	sD := (60 * 60 * 24)
	sH := (60 * 60)

	days := tSeconds / sD
	rSeconds := tSeconds % sD

	hours := rSeconds / sH
	rSeconds %= sH

	minutes := rSeconds / 60
	seconds := rSeconds % 60

	lT := fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
	if hours == 0 && days == 0 && minutes == 0 {
		lT = fmt.Sprintf("%d seconds", seconds)
	}
	if hours > 0 {
		lT = fmt.Sprintf("%d hours %s", hours, lT)
	}
	if days > 0 {
		lT = fmt.Sprintf("%d days %d hours", days, hours)
	}

	return widget.NewLabel(lT)
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
		//m := widget.NewMultiLineEntry()
		//m.SetText(note.Note)
		//m.Disable()
		//m.Wrapping = fyne.TextWrapWord
		//mn := canvas.NewText(note.Note, colornames.Blueviolet)
		m := widget.NewRichText(
			&widget.TextSegment{Text: note.Note, Style: widget.RichTextStyleParagraph},
		)
		m.Wrapping = fyne.TextWrapWord
		v.Add(m)
	}

	return v
}

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

func (td *TODO) LoadTODOTasks() {
	tasks, err := td.DB.AllTODOTasks()
	td.TODOTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) LoadCriticalTasks() {
	tasks, err := td.DB.AllOtherTabTasks("Critical")
	td.CriticalTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) LoadVeryHighTasks() {
	tasks, err := td.DB.AllOtherTabTasks("Very High")
	td.CriticalTasks = tasks
	if err != nil {
		log.Println(err)
	}

}
func (td *TODO) LoadHighTasks() {
	tasks, err := td.DB.AllOtherTabTasks("High")
	td.CriticalTasks = tasks
	if err != nil {
		log.Println(err)
	}

}
func (td *TODO) LoadMediumTasks() {
	tasks, err := td.DB.AllOtherTabTasks("Medium")
	td.CriticalTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) LoadLowTasks() {
	tasks, err := td.DB.AllOtherTabTasks("Low")
	td.CriticalTasks = tasks
	if err != nil {
		log.Println(err)
	}

}

func (td *TODO) getGenericStatusField(id int64, curPos int64, origTab string) *CustomSelect {
	s := NewCustomSelect(taskStatusColors, taskStatus, func(value string) {
		// Stop any existing timer
		previousStatus, _ := td.DB.GetStatusByTaskID(id)

		log.Println("Previous Status: ", previousStatus, "New Status: ", value, " for taskId: ", id)

		if previousStatus == "In Progress" {
			log.Println("Previous Status was 'In Progress'; We need to stop the previous timer; if any.")
			activeTimerID, err := td.DB.GetActiveTimerByTaskId(id)
			if err != nil {
				log.Println(err)
			}
			log.Println("Active Timer ID: ", activeTimerID)
			if activeTimerID != 0 {
				td.DB.StopTimer(activeTimerID)
			}
		}

		if value == "In Progress" {
			log.Println("In Progress...")
			// Need to stop any existing timer and switch their status automatically
			if td.UIElements.InProgressTimerId != 0 {
				// ...
			}

			timer, err := td.DB.StartTimer(id)
			if err != nil {
				log.Println(err)
			}

			td.UIElements.InProgressTimerTaskId = timer.TaskID
			td.UIElements.InProgressTimerId = timer.ID

		}

		td.DB.UpdateStatus(id, value)

		if value == "Done" {
			// We need to unshift the position
			log.Println("Shifting task id: ", id, " with position ", curPos)
			td.DB.ShiftPosition(id, curPos, "TODO")

			// Do as if we switch the Tab and reload everything
			td.refreshStatusTab(origTab)
		}

	})

	return s
}

func (td *TODO) getGenericTaskTable(origTab string) *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(td.CriticalTasks), len(TODODisplayColumns) // Column numbers
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(" . "))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			taskRow := td.CriticalTasks[i.Row]
			id := taskRow.ID

			colName := TODODisplayColumns[i.Col]

			switch colName {
			case "Position":
				curPos := taskRow.Position

				pc := td.getActionButtonsContainer(id, curPos, taskRow.Title)

				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					pc,
				}
			case "Title":
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(taskRow.Title),
				}
			case "Status":
				// Status
				sSel := td.getGenericStatusField(id, taskRow.Position, origTab)
				sSel.SetSelected(taskRow.Status)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					sSel,
				}
			case "Priority":
				pSel := td.getPriorityField(id)
				pSel.SetSelected(taskRow.Priority)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					pSel,
				}
			case "UpdatedAt":
				// updated_at
				//tUA, _ := time.ParseDuration(taskRow.UpdatedAt.Local().GoString())
				//dUA := time.Since(taskRow.UpdatedAt)
				uAl := td.getUpdatedAtField(taskRow.UpdatedAt)
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					uAl,
				}
			default:
				// Default is empty
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(""),
				}
			}

		})

	for i := 0; i < len(TODOColumnsSize); i++ {
		t.SetColumnWidth(i, TODOColumnsSize[i])
	}

	t.OnSelected = func(id widget.TableCellID) {
		log.Println("Cell was selected: ", id)
		// Which translate to:
		taskRow := td.CriticalTasks[id.Row]
		colName := TODODisplayColumns[id.Col]
		log.Println("Here is the row involved: ", taskRow, " with column name: ", colName)
		if colName == "Title" {
			log.Println("Value in cell: ", taskRow.Title)
			entryTitle := widget.NewEntry()
			fTitle := &widget.FormItem{
				Text:   "Title",
				Widget: entryTitle,
			}
			entryTitle.SetText(taskRow.Title)

			dialog.ShowForm("Edit task: "+taskRow.Title, "Save", "Cancel",
				[]*widget.FormItem{fTitle},
				func(submited bool) {
					if submited {
						log.Println("Save was press: ", entryTitle.Text)
						log.Println("Let's save ", entryTitle.Text, " to task id: ", taskRow.ID)
						// Let's save this
						taskRow.Title = entryTitle.Text
						td.DB.UpdateTitle(taskRow.ID, entryTitle.Text)
						td.LoadCriticalTasks()
						td.UIElements.CriticalTaskTable.Refresh()
					}
				}, td.MainWindow)
		}
	}

	return t
}
