package solution

import (
	"errors"
	"grader/pkg/server/user"
	"time"
)

type Solution struct {
	ID        int
	User      *user.Claims
	TaskID    int
	File      *File
	Result    *Result
	Status    string
	CreatedAt time.Time
}

type File struct {
	FileName string `json:"fileName "`
	File     []byte `json:"file"`
}

type Result struct {
	Pass bool   `json:"pass"`
	Text string `json:"text"`
}

var (
	ErrNoSolution = errors.New("No solution found")
)
