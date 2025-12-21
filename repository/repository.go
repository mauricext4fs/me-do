package repository

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertTask(a Tasks) (*Tasks, error)
	InsertPosition(p Positions) (*Positions, error)
	PushPosition() error
	ShiftPosition(id int64, curPos int64, label string) error
	UpPosition(id int64, curPos int64, label string) error
	DownPosition(id int64, curPos int64, label string) error
	AllTODOTasks() ([]Tasks, error)
	AllDoneTasks() ([]Tasks, error)
	AllOtherTabTasks(tabname string) ([]Tasks, error)
	SearchTODOTasks(searchText string) ([]Tasks, error)
	GetTaskByID(id int) (*Tasks, error)
	UpdateTask(id int64, updated Tasks) error
	UpdateStatus(id int64, status string) error
	GetStatusByTaskID(id int64) (string, error)
	UpdatePriority(id int64, status string) error
	UpdateTitle(id int64, title string) error
	DeleteTask(id int64) error
	AddNote(taskId int64, note string) error
	GetNotes(taskId int64) ([]Notes, error)
	AddFileToNote(noteId int64, filename string, filetype string) (int64, error)
	StartTimer(taskId int64) (*Timers, error)
	StopTimer(id int64) error
	GetActiveTimerByTaskId(id int64) (int64, error)
}

type Tasks struct {
	ID        int64     `json:"id"`
	Position  int64     `json:"position"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int64     `json:"updated_by"`
}

type Timers struct {
	ID             int64     `json:"id"`
	TaskID         int64     `json:"task_id"`
	StartTimestamp time.Time `json:"created_at"`
	EndTimestamp   time.Time `json:"updated_at"`
	CreatedBy      int64     `json:"created_by"`
}

type Notes struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int64     `json:"updated_by"`
}

type Positions struct {
	ID       int64  `json:"id"`
	TaskID   int64  `json:"task_id"`
	Position int64  `json:"position"`
	Label    string `json:"label"`
}

type Count struct {
	Count int64 `json:"id"`
}

type TaskLabel struct {
	ID       int64  `json:"id"`
	Position int64  `json:"position"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

func (t *Tasks) GetValueByName(name string) (value string) {
	switch name {
	case "ID":
		return strconv.FormatInt(t.ID, 10)
	case "Position":
		return strconv.FormatInt(t.ID, 10)
	case "Title":
		return fmt.Sprintf("%s", t.Title)
	case "Status":
		return fmt.Sprintf("%s", t.Status)
	case "Priority":
		return fmt.Sprintf("%s", t.Priority)
	default:
		return ""
	}

}
