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
	StartActivity(a Activities) (*Activities, error)
	AllActivities() ([]Activities, error)
	AllActivityType() ([]ActivityType, error)
	GetActivityByID(id int) (*Activities, error)
	UpdateActivity(id int64, updated Activities) error
	DeleteActivity(id int64) error
}

type Activities struct {
	ID             int64     `json:"id"`
	ActivityType   int64     `json:"type"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
}

type Count struct {
	Count int64 `json:"id"`
}

type ActivityType struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}
