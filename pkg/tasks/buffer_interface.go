package tasks

import (
	"fmt"
)

// Task is a struct that defines the structure of a task
type Task struct {
	Id          int    `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Uuid        string `json:"uuid"`

	isDone bool
}

// Buffer is an interface that defines the methods that a buffer should implement
type Buffer interface {
	SupportsAutoId() bool
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

type ValidateError struct {
	Message string
	Field   string
}

func (w *WriteError) Error() string {
	return fmt.Sprintf("[WriteError]: %s", w.Message)
}

func (r *ReadError) Error() string {
	return fmt.Sprintf("[ReadError]: %s", r.Message)
}

func (v *ValidateError) Error() string {
	return fmt.Sprintf("[ValidateError]: %s", v.Message)
}
