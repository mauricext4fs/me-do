package repository

import (
	"errors"
	"time"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertTask(a Tasks) (*Tasks, error)
	AllTasks() ([]Tasks, error)
	GetTaskByID(id int) (*Tasks, error)
	UpdateTask(id int64, updated Tasks) error
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
