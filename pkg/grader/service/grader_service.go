package service

import (
	"fmt"
	"grader/pkg/grader"
	"grader/pkg/grader/repo"
	"grader/pkg/server/solution"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type GraderServiceInterface interface {
	GradeFile(*solution.File) (*solution.Result, error)
	GraderLogin(string) error
	GetToken() (string, error)
}

type GraderService struct {
	Config     *grader.Config
	GraderRepo *repo.GraderRepo
}

func NewGraderService(config *grader.Config, graderRepo *repo.GraderRepo) *GraderService {
	return &GraderService{
		Config:     config,
		GraderRepo: graderRepo,
	}
}

func (s *GraderService) GraderLogin(token string) error {
	err := s.GraderRepo.SaveToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *GraderService) GetToken() (string, error) {
	t, err := s.GraderRepo.GetToken()
	if err != nil {
		return "", err
	}

	return t, err
}

func (s *GraderService) GradeFile(file *solution.File) (*solution.Result, error) {
	var fileName string

	for _, f := range s.Config.Files {
		if f.FileName == file.FileName {
			fileName = file.FileName
			break
		}
	}

	if fileName == "" {
		return nil, fmt.Errorf("empty filename")
	}

	result := &solution.Result{}

	contentBytes := []byte(file.File)

	tempDir, err := os.MkdirTemp("", "tempDir")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, fileName)

	err = os.WriteFile(tempFilePath, contentBytes, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to write file content to temporary file: %w", err)
	}

	timeoutSeconds := int64(1 * time.Minute / time.Second)
	stopTimeout := strconv.FormatInt(timeoutSeconds, 10)

	filePath := fmt.Sprintf("%s:/grader/solutionFiles/%s", tempFilePath, fileName)

	cmd := exec.Command(
		"docker",
		"run",
		"--user",
		"1000",
		"--network",
		"none",
		"--rm",
		"--name",
		"runXXXXXXXXXXXX",
		"--stop-timeout",
		stopTimeout,
		"-v",
		filePath,
		s.Config.GraderPayload.Container,
		"partId",
		s.Config.GraderPayload.PartID,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Pass = false
		result.Text = string(output)

		return result, nil
	}

	result.Pass = true
	result.Text = "Поздравляем! Вы успешно сделали задание"

	return result, nil
}
