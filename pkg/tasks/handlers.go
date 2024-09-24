package tasks

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type TaskOptions struct {
	Name        string
	Description string
	DueDate     time.Time
}

func InsertTask(
	opts TaskOptions,
	buff Buffer,
) (Task, error) {
	var latestId int

	// Get latest task from buffer
	latestTask, err := buff.GetLatest()

	if err != nil {
		var readErr *ReadError
		if errors.As(err, &readErr) {
			// If buffer is empty, set latest task id to 0
			latestId = 0
		} else {
			fmt.Println("Error getting latest task")
			return Task{}, err
		}
	} else {
		latestId = latestTask.Id
	}

	newId := latestId + 1
	newUuid, err := uuid.NewRandom()

	if err != nil {
		return Task{}, err
	}

	task := Task{
		Id:          newId,
		Name:        opts.Name,
		Description: opts.Description,
		DueDate:     opts.DueDate.Format(time.RFC3339),
		isDone:      false,
		Uuid:        newUuid.String(),
	}

	// Write task to buffer
	_, err = buff.Write(task)

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

func RemoveTaskById(
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

	// Update task in buffer
	_, err = buff.Update(id, task)

	if err != nil {
		return err
	}

	return nil
}
