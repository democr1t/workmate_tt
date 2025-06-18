package http

import (
	"net/http"
	"workmate_tt/internal/application"
)

func NewRouter(service *application.TaskService) http.Handler {
	handler := NewTaskHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/create", handler.CreateTask)
	mux.HandleFunc("GET /tasks", handler.GetTask)
	mux.HandleFunc("GET /tasks/all", handler.GetAllTasks)
	mux.HandleFunc("DELETE /tasks", handler.DeleteTask)

	return mux
}
