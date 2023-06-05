package user

import "errors"

type User struct {
	ID       string
	Username string
	Password string `json:"-"`
	Admin    bool
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

var (
	ErrNoUser             = errors.New("user not found")
	ErrBadPassword        = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
