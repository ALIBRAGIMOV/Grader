package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"grader/pkg/queue"
	"grader/pkg/server/session"
	"grader/pkg/server/solution"
	"grader/pkg/server/solution/service"
	"grader/pkg/utils"
	"io"
	"net/http"
)

type SolutionHandler struct {
	SolutionService service.SolutionServiceInterface
	RabbitChan      *amqp.Channel
}

func (h *SolutionHandler) SolutionResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		utils.GetLogger(ctx).Error("Error retrieving file", "not webhook")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	s := &solution.Solution{}

	err := utils.DecodeJSONHandler(w, r, &s)
	if err != nil {
		utils.GetLogger(ctx).Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.SolutionService.UpdateSolution(s)
	if err != nil {
		utils.GetLogger(ctx).Error("Error update solution", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SolutionHandler) UploadSolution(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.FormValue("id")

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.GetLogger(ctx).Error("Error retrieving file", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.GetLogger(ctx).Error("Error reading file", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sess, _ := session.SessionFromContext(ctx)
	if err != nil {
		utils.GetLogger(ctx).Error("Bad session", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s, err := h.SolutionService.UploadSolution(taskID, sess, fileBytes, fileHeader)
	if err != nil {
		utils.GetLogger(ctx).Error("Error uploading solution", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&s)
	if err != nil {
		utils.GetLogger(ctx).Error("Error marshal solution", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.RabbitChan.Publish(
		"",
		queue.SolutionQueueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})

	if err != nil {
		utils.GetLogger(ctx).Error("Error publish solution to queue", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/tasks/%s/solutions/%d", taskID, s.ID)
	http.Redirect(w, r, url, http.StatusFound)
}
