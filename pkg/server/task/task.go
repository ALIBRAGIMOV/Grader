package task

import (
	"errors"
	"time"
)

type Task struct {
	ID          int
	Name        string
	Description string
	Admins      []int
	CreatedAt   time.Time
}

var ErrNoTask = errors.New("task not found")
