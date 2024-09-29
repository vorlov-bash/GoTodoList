package tasks

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type InsertTaskOptions struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

func (ito *InsertTaskOptions) Validate() []error {
	var errs []error

	if ito.Name == "" {
		errs = append(errs, &ValidateError{Message: "Name is mandatory", Field: "Name"})
	}

	if ito.DueDate.IsZero() {
		errs = append(errs, &ValidateError{Message: "Due date is mandatory", Field: "DueDate"})
	}

	return errs
}

func InsertTask(
	name string,
	description string,
	dueDate time.Time,
	buff Buffer,
) (Task, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return Task{}, err
	}

	var task Task
	if !buff.SupportsAutoId() {
		var latestId int
		// Get latest task from buffer
		latestTask, err := buff.GetLatest()

		if err != nil {
			var readErr *ReadError
			if errors.As(err, &readErr) {
				// If buffer is empty, set latest task id to 0
				latestId = 0
			} else {
				return Task{}, err
			}
		} else {
			latestId = latestTask.Id
		}

		newId := latestId + 1

		task = Task{
			Id:          newId,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
			Name:        name,
			Description: description,
			DueDate:     dueDate.Format(time.RFC3339),
			isDone:      false,
			Uuid:        newUuid.String(),
		}
	} else {
		task = Task{
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
			Name:        name,
			Description: description,
			DueDate:     dueDate.Format(time.RFC3339),
			isDone:      false,
			Uuid:        newUuid.String(),
		}
	}

	task, err = buff.Write(task)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func GetTaskById(
	id int,
	buff Buffer,
) (Task, error) {
	// Get task from buffer
	task, err := buff.Get(id)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func GetAllTasks(
	buff Buffer,
) ([]Task, error) {
	// Get all tasks from buffer
	tasks, err := buff.GetAll()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTaskById(
	id int,
	buff Buffer,
) error {
	// Remove task from buffer
	return buff.Remove(id)
}

func MarkAsDone(
	id int,
	buff Buffer,
) error {
	// Get task from buffer
	task, err := buff.Get(id)

	if err != nil {
		return err
	}

	task.isDone = true
	task.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update task in buffer
	_, err = buff.Update(id, task)

	if err != nil {
		return err
	}

	return nil
}
