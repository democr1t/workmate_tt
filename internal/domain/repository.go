package domain

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	GetAll() ([]*Task, error)
	Update(task *Task) error
	Delete(id string) error
}
