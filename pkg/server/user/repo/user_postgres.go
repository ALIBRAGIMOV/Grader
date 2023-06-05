package repo

import (
	"database/sql"
	"grader/pkg/server/user"
	"strconv"
)

type Pgx struct {
	DB *sql.DB
}

type UserRepoInterface interface {
	Auth(string) (*user.User, error)
	Create(string, string) (*user.User, error)
	UserByID(string) (*user.User, error)
}

func NewPgxRepo(db *sql.DB) *Pgx {
	return &Pgx{
		DB: db,
	}
}

func (repo *Pgx) Auth(username string) (*user.User, error) {
	var userID int64
	u := &user.User{}

	row := repo.DB.QueryRow(`
        SELECT id, username, password FROM users
        WHERE username = $1
    `, username)

	err := row.Scan(&userID, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNoUser
		}
		return nil, err
	}

	u.ID = strconv.FormatInt(userID, 10)

	return u, nil
}

func (repo *Pgx) Create(username, password string) (*user.User, error) {
	var lastInsertId int64

	u := &user.User{
		ID:       "",
		Username: username,
		Password: password,
	}

	row := repo.DB.QueryRow(`
			INSERT INTO users ("username", "password")
			VALUES ($1, $2)
			RETURNING id
		`, username, password)

	err := row.Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}

	u.ID = strconv.FormatInt(lastInsertId, 10)

	return u, nil
}

func (repo *Pgx) UserByID(userID string) (*user.User, error) {
	u := &user.User{}
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}

	row := repo.DB.QueryRow(`
        SELECT id, username, password, admin
        FROM users
        WHERE id = $1
    `, id)

	err = row.Scan(&id, &u.Username, &u.Password, &u.Admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNoUser
		}
		return nil, err
	}

	u.ID = strconv.FormatInt(id, 10)

	return u, nil
}

//TODO added user list method and with query params limit offset
