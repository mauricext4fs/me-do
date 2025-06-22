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
	StartTask(a Tasks) (*Tasks, error)
	AllTasks() ([]Tasks, error)
	GetTaskByID(id int) (*Tasks, error)
	UpdateTask(id int64, updated Tasks) error
	DeleteTask(id int64) error
}

type Tasks struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
}

type Count struct {
	Count int64 `json:"id"`
}

type TaskLabel struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}
