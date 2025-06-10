package infrastructure

import (
	"sync"
	"workmate_tt/internal/domain"
)

type MemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *MemoryTaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) GetByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}
	return task, nil
}

func (r *MemoryTaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return domain.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *MemoryTaskRepository) GetAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}

	return tasks, nil
}
