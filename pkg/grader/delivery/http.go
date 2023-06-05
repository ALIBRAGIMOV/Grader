package delivery

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"grader/pkg/grader/service"
	"grader/pkg/server/solution"
	"grader/pkg/utils"
	"io"
	"net/http"
	"net/url"
)

type GraderHandler struct {
	GraderService service.GraderServiceInterface
	Logger        *zap.Logger
	FormData      url.Values
}

func (h *GraderHandler) GradeSolution(w http.ResponseWriter, r *http.Request) {
	s := &solution.Solution{}
	client := &http.Client{}

	err := utils.DecodeJSONHandler(w, r, &s)
	if err != nil {
		h.Logger.Error("Failed to decode JSON", zap.Error(err))
		return
	}

	result, err := h.GraderService.GradeFile(s.File)
	if err != nil {
		h.Logger.Error("Failed to grade file", zap.Error(err))
		http.Error(w, "error grade file", http.StatusInternalServerError)
		return
	}

	s.Result = result
	s.Status = "completed"

	resultData, err := json.Marshal(s)
	if err != nil {
		h.Logger.Error("Failed to marshal solution", zap.Error(err))
		http.Error(w, "Failed to marshal solution:", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/webhook/solution/result", bytes.NewBuffer(resultData))
	if err != nil {
		h.Logger.Error("Failed to send webhook request:", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.GraderService.GetToken()
	if err != nil {
		h.Logger.Error("Failed to get token:", zap.Error(err))
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
	}

	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		h.Logger.Error("Failed to send webhook request:", zap.Error(err))
		http.Error(w, "Failed to send webhook request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
}

func (h *GraderHandler) GraderLogin() {
	resp, err := http.PostForm("http://localhost:3000/api/v1/user/login", h.FormData)
	if err != nil {
		h.Logger.Error("Could not login to server service", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.Logger.Error("Failed to read response body", zap.Error(err))
		return
	}

	if resp.StatusCode == http.StatusOK {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "token" {
				err = h.GraderService.GraderLogin(cookie.Value)
				if err != nil {
					h.Logger.Error("Failed login with token", zap.String("token", cookie.Value))
					return
				}

				return
			}
		}
	} else {
		h.Logger.Error("Failed to login", zap.String("responseBody", string(body)))
		return
	}
}
