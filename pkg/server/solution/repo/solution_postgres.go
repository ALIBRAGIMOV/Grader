package repo

import (
	"database/sql"
	"encoding/json"
	"grader/pkg/server/solution"
)

type Pgx struct {
	DB *sql.DB
}

type SolutionRepoInterface interface {
	Add(*solution.Solution) (*solution.Solution, error)
	Update(*solution.Solution) error
	GetListByTaskID(int) ([]*solution.Solution, error)
	GetByID(int) (*solution.Solution, error)
	List() ([]*solution.Solution, error)
}

func NewPgxRepo(db *sql.DB) *Pgx {
	return &Pgx{
		DB: db,
	}
}

func (repo *Pgx) Add(s *solution.Solution) (*solution.Solution, error) {
	var lastInsertId int64
	res := &solution.Solution{
		User:      s.User,
		TaskID:    s.TaskID,
		File:      s.File,
		Result:    s.Result,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
	}

	userJson, err := json.Marshal(s.User)
	if err != nil {
		return nil, err
	}
	resultJson, err := json.Marshal(s.Result)
	if err != nil {
		return nil, err
	}
	fileJson, err := json.Marshal(s.File)
	if err != nil {
		return nil, err
	}

	row := repo.DB.QueryRow(`
		INSERT INTO solutions (user_data, task_id, file, result, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, userJson, s.TaskID, fileJson, resultJson, s.Status, s.CreatedAt)

	err = row.Scan(
		&lastInsertId,
	)
	if err != nil {
		return nil, err
	}

	res.ID = int(lastInsertId)

	return res, nil
}

func (repo *Pgx) Update(s *solution.Solution) error {
	userJson, err := json.Marshal(s.User)
	if err != nil {
		return err
	}
	resultJson, err := json.Marshal(s.Result)
	if err != nil {
		return err
	}
	fileJson, err := json.Marshal(s.File)
	if err != nil {
		return err
	}

	_, err = repo.DB.Exec(`
		UPDATE solutions 
		SET user_data = $1, task_id = $2, file = $3, result = $4, status = $5, created_at = $6
		WHERE id = $7
	`, userJson, s.TaskID, fileJson, resultJson, s.Status, s.CreatedAt, s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Pgx) List() ([]*solution.Solution, error) {
	//TODO added with query params limit offset
	rows, err := repo.DB.Query(`
		SELECT id, user_data, task_id, file, result, status, created_at
		FROM solutions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solutions []*solution.Solution

	for rows.Next() {
		var s solution.Solution
		var userJSON []byte
		var resultJSON []byte
		var fileJson []byte

		err = rows.Scan(
			&s.ID,
			&userJSON,
			&s.TaskID,
			&fileJson,
			&resultJSON,
			&s.Status,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(userJSON, &s.User)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resultJSON, &s.Result)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(fileJson, &s.File)
		if err != nil {
			return nil, err
		}

		solutions = append(solutions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return solutions, nil
}

func (repo *Pgx) GetListByTaskID(taskID int) ([]*solution.Solution, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_data, task_id, file, result, status, created_at
		FROM solutions
		WHERE task_id = $1
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solutions []*solution.Solution

	for rows.Next() {
		var s solution.Solution
		var userJSON []byte
		var resultJSON []byte
		var fileJson []byte

		err = rows.Scan(
			&s.ID,
			&userJSON,
			&s.TaskID,
			&fileJson,
			&resultJSON,
			&s.Status,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(userJSON, &s.User)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resultJSON, &s.Result)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(fileJson, &s.File)
		if err != nil {
			return nil, err
		}

		solutions = append(solutions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return solutions, nil
}

func (repo *Pgx) GetByID(id int) (*solution.Solution, error) {
	row := repo.DB.QueryRow(`
		SELECT id, user_data, task_id, file, result, status, created_at
		FROM solutions
		WHERE id = $1
	`, id)

	var s solution.Solution
	var userJSON []byte
	var resultJSON []byte
	var fileJson []byte

	err := row.Scan(
		&s.ID,
		&userJSON,
		&s.TaskID,
		&fileJson,
		&resultJSON,
		&s.Status,
		&s.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, solution.ErrNoSolution
		}
		return nil, err
	}

	err = json.Unmarshal(userJSON, &s.User)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultJSON, &s.Result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileJson, &s.File)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
