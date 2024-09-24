package tasks

import "fmt"

// Task is a struct that defines the structure of a task
type Task struct {
	Id          int
	Name        string
	Description string
	DueDate     string
	Uuid        string

	isDone bool
}

// Buffer is an interface that defines the methods that a buffer should implement
type Buffer interface {
	Write(data Task) (Task, error)
	WriteBatch(data []Task) ([]Task, error)
	Remove(id int) error
	RemoveBatch(ids []int) error
	Update(id int, data Task) (Task, error)
	Get(id int) (Task, error)
	GetAll() ([]Task, error)
	GetLatest() (Task, error)
}

// WriteError is occurred when a write operation fails (write, writeBatch, update, remove, removeBatch)
type WriteError struct {
	Message string
}

// ReadError is occurred when a read operation fails (get, getAll)
type ReadError struct {
	Message string
}

func (w *WriteError) Error() string {
	return fmt.Sprintf("[WriteError]: %s", w.Message)
}

func (r *ReadError) Error() string {
	return fmt.Sprintf("[ReadError]: %s", r.Message)
}
