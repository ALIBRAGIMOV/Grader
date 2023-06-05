package service

import (
	"grader/pkg/server/session"
	"grader/pkg/server/solution"
	"grader/pkg/server/solution/repo"
	"mime/multipart"
	"strconv"
	"time"
)

type SolutionServiceInterface interface {
	UploadSolution(string, *session.Session, []byte, *multipart.FileHeader) (*solution.Solution, error)
	GetSolutionsByTaskID(string, string, bool) ([]*solution.Solution, error)
	GetSolutionByID(string) (*solution.Solution, error)
	GetSolutionsByUserName(string) ([]*solution.Solution, error)
	UpdateSolution(*solution.Solution) error
}

type SolutionService struct {
	SolutionRepoPQ repo.SolutionRepoInterface
}

func NewSolutionService(pgx repo.SolutionRepoInterface) *SolutionService {
	return &SolutionService{
		SolutionRepoPQ: pgx,
	}
}

func (h *SolutionService) GetSolutionsByUserName(user string) ([]*solution.Solution, error) {
	var filteredByUser []*solution.Solution

	solutions, err := h.SolutionRepoPQ.List()
	if err != nil {
		return nil, err
	}

	for _, s := range solutions {
		if s.User.Username != "" && s.User.Username == user {
			filteredByUser = append(filteredByUser, s)
		}
	}

	return filteredByUser, nil
}

func (h *SolutionService) UploadSolution(taskID string, sess *session.Session, file []byte, fileHeader *multipart.FileHeader) (*solution.Solution, error) {
	tID, err := strconv.Atoi(taskID)
	if err != nil {
		return nil, err
	}

	s := &solution.Solution{
		TaskID: tID,
		User:   sess.User,
		File: &solution.File{
			File:     file,
			FileName: fileHeader.Filename,
		},
		CreatedAt: time.Now(),
		Result: &solution.Result{
			Pass: false,
			Text: "У вас ошибка в задании",
		},
		Status: "pending",
	}

	s, err = h.SolutionRepoPQ.Add(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (h *SolutionService) UpdateSolution(s *solution.Solution) error {
	err := h.SolutionRepoPQ.Update(s)
	if err != nil {
		return err
	}

	return nil
}

func (h *SolutionService) GetSolutionsByTaskID(taskID string, uID string, isAdmin bool) ([]*solution.Solution, error) {
	var filteredByUser []*solution.Solution

	tID, err := strconv.Atoi(taskID)
	if err != nil {
		return nil, err
	}

	solutions, err := h.SolutionRepoPQ.GetListByTaskID(tID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		for _, s := range solutions {
			if s.User.ID != "" && s.User.ID == uID {
				filteredByUser = append(filteredByUser, s)
			}
		}

		return filteredByUser, nil
	}

	return solutions, nil
}

func (h *SolutionService) GetSolutionByID(solutionID string) (*solution.Solution, error) {
	sID, err := strconv.Atoi(solutionID)
	if err != nil {
		return nil, err
	}

	s, err := h.SolutionRepoPQ.GetByID(sID)
	if err != nil {
		return nil, err
	}

	return s, nil
}
