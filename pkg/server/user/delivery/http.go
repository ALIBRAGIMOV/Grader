package delivery

import (
	"go.uber.org/zap"
	"grader/pkg/server/user"
	"grader/pkg/server/user/service"
	"grader/pkg/utils"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	Tmpl        *template.Template
	UserService service.UserServiceInterface
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

type Task struct {
	TaskDescription string
	TaskDeadline    string
	TaskDetails     string
}

func (h *UserHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	task := Task{
		TaskDescription: "Sample Task",
		TaskDeadline:    "2023-05-31",
		TaskDetails:     "This is a sample task description.",
	}

	err := h.Tmpl.ExecuteTemplate(w, "index.html", task)
	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.UserService.Login(username, password)

	if err != nil {
		if err == user.ErrNoUser {
			utils.GetLogger(ctx).Error("User not found")
		} else {
			utils.GetLogger(ctx).Error("Error logging", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	utils.GetLogger(ctx).Info("User logged in", zap.String("username", username))

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Cookie")
	ctx := r.Context()

	if authHeader == "" {
		http.Error(w, "Authorization header not found", http.StatusUnauthorized)
		return
	}

	t := strings.TrimPrefix(authHeader, "token=")

	err := h.UserService.Logout(t)
	if err != nil {
		utils.GetLogger(ctx).Error("Error user logout", zap.Error(err))
		http.Error(w, `Logout error`, http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.UserService.Register(username, password)
	if err != nil {
		utils.GetLogger(ctx).Error("Error decoding request body", zap.Error(err))
		http.Error(w, "error token", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	utils.GetLogger(ctx).Info("User register in", zap.String("username", username))

	http.Redirect(w, r, "/tasks", http.StatusFound)
}
