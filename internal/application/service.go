package application

import (
	"workmate_tt/internal/domain"
	"workmate_tt/internal/infrastructure"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo   domain.TaskRepository
	workerPool *infrastructure.WorkerPool
	taskFunc   domain.TaskFunction
}

func NewTaskService(
	taskRepo domain.TaskRepository,
	maxWorkers int,
	taskFunc domain.TaskFunction,
) *TaskService {
	workerPool := infrastructure.NewWorkerPool(taskRepo, maxWorkers, taskFunc)
	return &TaskService{
		taskRepo:   taskRepo,
		workerPool: workerPool,
		taskFunc:   taskFunc,
	}
}

func (s *TaskService) CreateTask() (*TaskDTO, error) {
	taskID := uuid.New()
	task := domain.NewTask(taskID.String())

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	s.workerPool.Submit(task)
	return s.mapToDTO(task), nil
}

func (s *TaskService) GetTask(id string) (*TaskDTO, error) {
	task, err := s.taskRepo.GetByID(id)

	if err != nil {
		return nil, err
	}

	return s.mapToDTO(task), nil
}

func (s *TaskService) DeleteTask(id string) error {
	return s.taskRepo.Delete(id)
}

func (s *TaskService) GetAllTasks() ([]*TaskDTO, error) {
	tasks, err := s.taskRepo.GetAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]*TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		dtos = append(dtos, s.mapToDTO(task))
	}

	return dtos, nil
}

func (s *TaskService) mapToDTO(task *domain.Task) *TaskDTO {
	var errStr string

	if task.Error != nil {
		errStr = task.Error.Error()
	}

	return &TaskDTO{
		ID:          task.ID,
		Status:      string(task.Status),
		Result:      task.Result,
		Error:       errStr,
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Duration:    task.Duration(),
	}
}
