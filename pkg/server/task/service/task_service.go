package service

import (
	"fmt"
	"grader/pkg/server/solution"
	"grader/pkg/server/task"
	"grader/pkg/server/task/repo"
	"strconv"
)

type TaskServiceInterface interface {
	GetTaskList() ([]*task.Task, error)
	GetTaskByID(string) (*task.Task, error)
	CreateTask(string, string) error
	UpdateTask(string, string, string) error
	GetTasksByUserSolutions([]*solution.Solution) ([]*task.Task, error)
}

type TaskService struct {
	TaskRepoPQ repo.TaskRepoInterface
}

func NewTaskService(pgx repo.TaskRepoInterface) *TaskService {
	return &TaskService{
		TaskRepoPQ: pgx,
	}
}

func (h *TaskService) GetTasksByUserSolutions(solutions []*solution.Solution) ([]*task.Task, error) {
	var tasks []*task.Task

	taskIDs := make(map[int]bool)
	var filteredTaskIDs []string

	for _, s := range solutions {
		if _, ok := taskIDs[s.TaskID]; !ok {
			taskIDs[s.TaskID] = true
			tID := fmt.Sprintf("%d", s.TaskID)

			filteredTaskIDs = append(filteredTaskIDs, tID)
		}
	}

	for _, i := range filteredTaskIDs {
		t, err := h.GetTaskByID(i)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (h *TaskService) UpdateTask(name, description, taskID string) error {
	t, err := h.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	t.Name = name
	t.Description = description

	err = h.TaskRepoPQ.Update(t)
	if err != nil {
		return err
	}

	return nil
}

func (h *TaskService) GetTaskList() ([]*task.Task, error) {
	//TODO added with query params limit offset
	t, err := h.TaskRepoPQ.List(0, 15)
	if err != nil {
		return nil, err
	}

	return t, err
}

func (h *TaskService) GetTaskByID(taskID string) (*task.Task, error) {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return nil, err
	}

	t, err := h.TaskRepoPQ.Get(id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h *TaskService) CreateTask(name, description string) error {
	t := &task.Task{
		Name:        name,
		Description: description,
	}

	err := h.TaskRepoPQ.Add(t)
	if err != nil {
		return err
	}

	return nil
}
