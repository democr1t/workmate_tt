package infrastructure

import (
	"sync"
	"workmate_tt/internal/domain"
)

type WorkerPool struct {
	taskRepo   domain.TaskRepository
	maxWorkers int
	taskChan   chan *domain.Task
	taskFunc   domain.TaskFunction
	wg         sync.WaitGroup
}

func NewWorkerPool(
	taskRepo domain.TaskRepository,
	maxWorkers int,
	taskFunc domain.TaskFunction,
) *WorkerPool {
	pool := &WorkerPool{
		taskRepo:   taskRepo,
		maxWorkers: maxWorkers,
		taskChan:   make(chan *domain.Task),
		taskFunc:   taskFunc,
	}

	pool.start()
	return pool
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for task := range p.taskChan {
		task.Start()
		_ = p.taskRepo.Update(task)

		result, err := p.taskFunc()

		task.Complete(result, err)
		_ = p.taskRepo.Update(task)
	}
}

func (p *WorkerPool) start() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) Submit(task *domain.Task) {
	p.taskChan <- task
}

func (p *WorkerPool) Stop() {
	close(p.taskChan)
	p.wg.Wait()
}
