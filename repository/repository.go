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
	AllTODOTasks() ([]Tasks, error)
	GetTaskByID(id int) (*Tasks, error)
	UpdateTask(id int64, updated Tasks) error
	UpdateStatus(id int64, status string) error
	UpdatePosition(id int64, newPos int64) error
	DeleteTask(id int64) error
}

type Tasks struct {
	ID        int64     `json:"id"`
	Position  int64     `json:"position"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
