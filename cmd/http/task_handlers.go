package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"net/http"
	"strconv"
	"time"
)

func TaskRouter(buff tasks.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetTask(buff)(w, r)
		case http.MethodPost:
			CreateTask(buff)(w, r)
		case http.MethodDelete:
			DeleteTask(buff)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func GetTask(buff tasks.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawTaskId := r.URL.Query().Get("id")

		if rawTaskId == "" {
			ValidationError(w, "Task id is mandatory", nil)
			return
		}

		taskId, err := strconv.Atoi(rawTaskId)

		if err != nil {
			ValidationError(w, "Task id must be a number", nil)
			return
		}

		task, err := tasks.GetTaskById(taskId, buff)

		if err != nil {
			var err *tasks.ReadError
			if errors.As(err, &err) {
				NotFoundError(w, fmt.Sprintf("Task with id=%d is not found", taskId), nil)
				return
			}

			InternalServerError(w, "", map[string]interface{}{"errors": []error{err}})
			return
		}

		_ = WriteStructToJson(w, task)
	}
}

func CreateTask(buff tasks.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract task data from request
		var input map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&input)

		if err != nil {
			ValidationError(w, "Due date must be a valid date (2000-01-01)", nil)
			return
		}

		if err != nil {
			InternalServerError(w, "", map[string]interface{}{"errors": []error{err}})
			return
		}

		// Validate task data
		dueDate, err := time.Parse(time.DateOnly, input["dueDate"].(string))

		if err != nil {
			ValidationError(w, "Due date must be a valid date (2000-01-01)", nil)
			return
		}

		// Insert task into buffer
		newTask, err := tasks.InsertTask(input["name"].(string), input["description"].(string), dueDate, buff)
		if err != nil {
			InternalServerError(w, "", map[string]interface{}{"errors": []error{err}})
			return
		}

		_ = WriteStructToJson(w, newTask)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func DeleteTask(buff tasks.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawTaskId := r.URL.Query().Get("id")

		if rawTaskId == "" {
			ValidationError(w, "Task id is mandatory", nil)
		}

		taskId, err := strconv.Atoi(rawTaskId)

		if err != nil {
			ValidationError(w, "Task id must be a number", nil)
		}

		err = tasks.DeleteTaskById(taskId, buff)

		if err != nil {
			InternalServerError(w, "", map[string]interface{}{"errors": []error{err}})
			return
		}
	}
}
