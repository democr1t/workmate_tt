package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"workmate_tt/internal/application"
	"workmate_tt/internal/infrastructure"
	consumers "workmate_tt/internal/interfaces/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	taskRepo := infrastructure.NewMemoryTaskRepository()

	taskFunc := func() (interface{}, error) {

		time.Sleep(time.Duration(3) * time.Second)
		return fmt.Sprintf("Processed in %.2f seconds", 3.0), nil
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
