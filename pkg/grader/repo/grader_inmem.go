package repo

import "sync"

type GraderRepo struct {
	mu    *sync.Mutex
	token string
}

func NewGraderRepo() *GraderRepo {
	return &GraderRepo{
		mu:    &sync.Mutex{},
		token: "",
	}
}

func (repo *GraderRepo) SaveToken(token string) error {
	repo.mu.Lock()
	repo.token = token
	repo.mu.Unlock()

	return nil
}

func (repo *GraderRepo) GetToken() (string, error) {
	var t string
	repo.mu.Lock()
	t = repo.token
	repo.mu.Unlock()

	return t, nil
}
