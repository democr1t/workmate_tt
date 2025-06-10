package http

import (
	"encoding/json"
	"net/http"
	"time"
	"workmate_tt/internal/application"
)

type TaskHandler struct {
	service *application.TaskService
}

func NewTaskHandler(service *application.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTaskRequest структура запроса на создание задачи
type CreateTaskRequest struct {
	Params interface{} `json:"params"`
}

// TaskResponse структура ответа с информацией о задаче
type TaskResponse struct {
	ID          string        `json:"id"`
	Status      string        `json:"status"`
	Result      interface{}   `json:"result,omitempty"`
	Error       string        `json:"error,omitempty"`
	CreatedAt   string        `json:"created_at"`
	StartedAt   string        `json:"started_at,omitempty"`
	CompletedAt string        `json:"completed_at,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.mapToResponse(task))
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.mapToResponse(task))
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, h.mapToResponse(task))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) mapToResponse(dto *application.TaskDTO) TaskResponse {
	return TaskResponse{
		ID:          dto.ID,
		Status:      dto.Status,
		Result:      dto.Result,
		Error:       dto.Error,
		CreatedAt:   dto.CreatedAt.Format(time.RFC3339),
		StartedAt:   formatTime(dto.StartedAt),
		CompletedAt: formatTime(dto.CompletedAt),
		Duration:    dto.Duration,
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
