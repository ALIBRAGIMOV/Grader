package delivery

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"grader/pkg/server/session"
	"grader/pkg/server/solution"
	solutionService "grader/pkg/server/solution/service"
	"grader/pkg/server/task"
	"grader/pkg/server/task/service"
	"grader/pkg/server/user"
	userService "grader/pkg/server/user/service"
	"grader/pkg/utils"
	"html/template"
	"net/http"
)

type TaskHandler struct {
	Tmpl            *template.Template
	TaskService     service.TaskServiceInterface
	SolutionService solutionService.SolutionServiceInterface
	UserService     userService.UserServiceInterface
}

type TaskData struct {
	User      *user.Claims
	Task      *task.Task
	Solutions []*solution.Solution
	Solution  *solution.Solution
}

type TasksData struct {
	User  *user.Claims
	Tasks []*task.Task
}

func (h *TaskHandler) TaskCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !u.Admin {
		url := fmt.Sprintf("/tasks/user/%s", u.Username)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "task_create.html",
		struct {
			User *user.Claims
		}{
			User: sess.User,
		})

	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) TaskEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	t := &task.Task{}

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !u.Admin {
		url := fmt.Sprintf("/tasks/user/%s", u.Username)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	t, err = h.TaskService.GetTaskByID(taskID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get task by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "task_edit.html",
		struct {
			User *user.Claims
			Task *task.Task
			URL  string
		}{
			User: sess.User,
			Task: t,
			URL:  r.URL.String(),
		})

	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) TaskAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !u.Admin {
		utils.GetLogger(ctx).Error("User not admin")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.TaskService.CreateTask(name, description)
	if err != nil {
		utils.GetLogger(ctx).Error("error create task", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/tasks/admin/task/all")
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *TaskHandler) TaskUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	taskID := r.FormValue("id")

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !u.Admin {
		utils.GetLogger(ctx).Error("User not admin", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.TaskService.UpdateTask(name, description, taskID)
	if err != nil {
		utils.GetLogger(ctx).Error("error create task", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/tasks/admin/task/all")
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *TaskHandler) TaskList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := &TasksData{}

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data.User = sess.User

	if !u.Admin {
		url := fmt.Sprintf("/tasks/user/%s", u.Username)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	//TODO added with query params limit offset

	t, err := h.TaskService.GetTaskList()
	if err != nil {
		utils.GetLogger(ctx).Error("error get task list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data.Tasks = t

	err = h.Tmpl.ExecuteTemplate(w, "task_list.html", data)
	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) TasksByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uName := chi.URLParam(r, "user")
	data := &TasksData{}

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data.User = sess.User

	if sess.User.Username != uName {
		url := fmt.Sprintf("/tasks/user/%s", sess.User.Username)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	solutions, err := h.SolutionService.GetSolutionsByUserName(uName)
	if err != nil {
		utils.GetLogger(ctx).Error("error get solutions", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	tasks, err := h.TaskService.GetTasksByUserSolutions(solutions)
	if err != nil {
		utils.GetLogger(ctx).Error("error get tasks", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data.Tasks = tasks

	err = h.Tmpl.ExecuteTemplate(w, "tasks_by_user.html", data)
	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	solutionID := chi.URLParam(r, "solutionID")
	data := &TaskData{}
	s := &solution.Solution{}

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data.User = sess.User

	t, err := h.TaskService.GetTaskByID(taskID)
	if err != nil {
		if err == task.ErrNoTask {
			utils.GetLogger(ctx).Error("Task not found")
		} else {
			utils.GetLogger(ctx).Error("Error get task", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	data.Task = t

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	solutions, err := h.SolutionService.GetSolutionsByTaskID(taskID, sess.User.ID, u.Admin)
	if err != nil {
		if err == task.ErrNoTask {
			utils.GetLogger(ctx).Error("Solutions not found", zap.Error(err))
		} else {
			utils.GetLogger(ctx).Error("Error get solutions", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	data.Solutions = solutions

	if solutionID != "" {
		s, err = h.SolutionService.GetSolutionByID(solutionID)
		if sess.User.ID != s.User.ID {
			url := fmt.Sprintf("/tasks/%s", taskID)
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		if err != nil {
			if err == solution.ErrNoSolution {
				utils.GetLogger(ctx).Error("Solution not found", zap.Error(err))
			} else {
				utils.GetLogger(ctx).Error("Error get solution", zap.Error(err))
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		data.Solution = s
	}

	err = h.Tmpl.ExecuteTemplate(w, "task_by_id.html", data)
	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) TaskSolutions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := chi.URLParam(r, "id")
	data := &TaskData{}

	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("error get session from context", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u, err := h.UserService.UserByID(sess.User.ID)
	if err != nil {
		utils.GetLogger(ctx).Error("error get user by id", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !u.Admin {
		url := fmt.Sprintf("/tasks/user/%s", u.Username)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	data.User = sess.User

	t, err := h.TaskService.GetTaskByID(taskID)
	if err != nil {
		if err == task.ErrNoTask {
			utils.GetLogger(ctx).Error("Task not found")
		} else {
			utils.GetLogger(ctx).Error("Error get task", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	data.Task = t

	solutions, err := h.SolutionService.GetSolutionsByTaskID(taskID, sess.User.ID, u.Admin)
	if err != nil {
		if err == task.ErrNoTask {
			utils.GetLogger(ctx).Error("Solutions not found", zap.Error(err))
		} else {
			utils.GetLogger(ctx).Error("Error get solutions", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	data.Solutions = solutions

	err = h.Tmpl.ExecuteTemplate(w, "task_solutions.html", data)
	if err != nil {
		utils.GetLogger(ctx).Error("error execute template", zap.Error(err))
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}
