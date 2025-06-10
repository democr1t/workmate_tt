package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"workmate_tt/internal/application"
	"workmate_tt/internal/infrastructure"
	consumers "workmate_tt/internal/interfaces/http"
)

func main() {
	godotenv.Load()
	taskRepo := infrastructure.NewMemoryTaskRepository()

	taskFunc := func(params interface{}) (interface{}, error) {
		duration, ok := params.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid params type, expected number")
		}

		time.Sleep(time.Duration(duration) * time.Second)
		return fmt.Sprintf("Processed in %.2f seconds", duration), nil
	}

	maxWorkers := runtime.NumCPU()
	taskService := application.NewTaskService(taskRepo, maxWorkers, taskFunc)

	router := consumers.NewRouter(taskService)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
