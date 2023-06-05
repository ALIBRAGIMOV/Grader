package service

import (
	"grader/pkg/server/session"
	"grader/pkg/server/user"
	"grader/pkg/server/user/repo"
	"grader/pkg/utils"
)

type UserServiceInterface interface {
	Login(username, password string) (string, error)
	Logout(string) error
	Register(username, password string) (string, error)
	UserByID(string) (*user.User, error)
}

type UserService struct {
	UserRepoPQ  repo.UserRepoInterface
	SessionRepo session.Manager
}

func NewUserService(pgx repo.UserRepoInterface, sessions session.Manager) *UserService {
	return &UserService{
		UserRepoPQ:  pgx,
		SessionRepo: sessions,
	}
}

func (h *UserService) Login(username, password string) (string, error) {
	u, err := h.UserRepoPQ.Auth(username)
	if err != nil {
		return "", err
	}

	isValidPass := utils.HashPassValidator(password, u.Password)

	if !isValidPass {
		return "", user.ErrBadPassword
	}

	token, err := h.SessionRepo.Create(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *UserService) Logout(token string) error {
	err := h.SessionRepo.Delete(token)
	if err != nil {
		return err
	}

	return nil
}

func (h *UserService) Register(username, password string) (string, error) {
	hashPass, err := utils.HashPass(password)

	u, err := h.UserRepoPQ.Create(username, hashPass)
	if err != nil {
		return "", err
	}

	token, err := h.SessionRepo.Create(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *UserService) UserByID(uID string) (*user.User, error) {
	u, err := h.UserRepoPQ.UserByID(uID)
	if err != nil {
		return nil, err
	}

	return u, nil
}
