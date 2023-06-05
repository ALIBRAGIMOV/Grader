package session

import (
	"errors"
	"grader/pkg/server/user"
)

type Session struct {
	ID   string
	User *user.Claims
}

type sessKey string

var Key sessKey = "token"

var (
	ErrNoAuth = errors.New("No session found")
)
