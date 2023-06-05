package repo

import (
	"database/sql"
	"github.com/lib/pq"
	"grader/pkg/server/task"
)

type Pgx struct {
	DB *sql.DB
}

type TaskRepoInterface interface {
	Add(*task.Task) error
	Update(*task.Task) error
	List(int, int) ([]*task.Task, error)
	Get(int) (*task.Task, error)
}

func NewPgxRepo(db *sql.DB) *Pgx {
	return &Pgx{
		DB: db,
	}
}

func (repo *Pgx) Add(t *task.Task) error {
	var taskID int

	err := repo.DB.QueryRow(`
		INSERT INTO tasks (name, description, admins, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id;
	`, t.Name, t.Description, pq.Array(t.Admins)).Scan(&taskID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Pgx) Update(t *task.Task) error {
	_, err := repo.DB.Exec(`
		UPDATE tasks 
		SET name = $1, description = $2, admins = $3, created_at = $4
		WHERE id = $5;
	`, t.Name, t.Description, pq.Array(t.Admins), t.CreatedAt, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Pgx) List(limit, offset int) ([]*task.Task, error) {
	//TODO add limit offset
	rows, err := repo.DB.Query(`
		SELECT id, name, description, admins, created_at
		FROM tasks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task

	for rows.Next() {
		var t task.Task
		var admins pq.Int64Array

		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&admins,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		t.Admins = make([]int, len(admins))
		for i, admin := range admins {
			t.Admins[i] = int(admin)
		}

		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *Pgx) Get(taskID int) (*task.Task, error) {
	t := &task.Task{}

	row := repo.DB.QueryRow(`
		SELECT id, name, description, admins, created_at
		FROM tasks
		WHERE id = $1
	`, taskID)

	var admins pq.Int64Array

	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&admins,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, task.ErrNoTask
		}
		return nil, err
	}

	t.Admins = make([]int, len(admins))
	for i, admin := range admins {
		t.Admins[i] = int(admin)
	}

	return t, nil
}
