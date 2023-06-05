package delivery

import (
	"bytes"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/http"
)

type QueueHandler struct {
	Client *http.Client
	Logger *zap.Logger
}

func (h *QueueHandler) SolutionWorker(solutions <-chan amqp.Delivery) {
	for s := range solutions {
		h.solutionHandler(s)
	}
}

func (h *QueueHandler) solutionHandler(s amqp.Delivery) {
	defer func() {
		if err := recover(); err != nil {
			h.Logger.Error("Recovered", zap.Any("panic", err))
		}
	}()

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/grader/grade", bytes.NewBuffer(s.Body))
	if err != nil {
		h.Logger.Error("Failed to create HTTP request", zap.Error(err))
		s.Ack(false)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		h.Logger.Error("Failed to make HTTP request", zap.Error(err))
		s.Ack(false)
		return
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.Logger.Error("Unexpected response status", zap.Int("status", resp.StatusCode))
		s.Ack(false)
		return
	}

	s.Ack(false)
}
