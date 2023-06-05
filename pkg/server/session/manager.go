package session

import (
	"grader/pkg/server/user"
)

type Manager interface {
	Create(user *user.User) (string, error)
	Delete(string) error
	Check(string) (*Session, error)
	Get(*Session) (string, error)
}
