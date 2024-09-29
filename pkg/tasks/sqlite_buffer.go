package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Sqlite3Buffer struct {
	db *sql.DB
}

func NewSqlite3Buffer(fileName string) (*Sqlite3Buffer, error) {
	err := os.MkdirAll("tmp", os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("[NewSqlite3Buffer]: cannot create directory: %v", err)
	}

	db, err := sql.Open("sqlite3", fileName)

	if err != nil {
		return nil, fmt.Errorf("[NewSqlite3Buffer]: cannot open database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
		    			id INTEGER PRIMARY KEY AUTOINCREMENT,
		    			createdAt TEXT,
		    			updatedAt TEXT,
		    			name TEXT NOT NULL,
		    			description TEXT,
		    			dueDate TEXT NOT NULL,
		    			uuid TEXT NOT NULL,
		    			isDone BOOLEAN NOT NULL
		)
`)
	return &Sqlite3Buffer{db: db}, nil
}

func (sb *Sqlite3Buffer) SupportsAutoId() bool {
	return true
}

func (sb *Sqlite3Buffer) Write(task Task) (Task, error) {
	result, err := sb.db.Exec(
		"INSERT INTO tasks (createdAt, updatedAt, name, description, dueDate, isDone, uuid) VALUES (?, ?, ?, ?, ?, ?, ?)",
		task.CreatedAt, task.UpdatedAt, task.Name, task.Description, task.DueDate, task.isDone, task.Uuid,
	)
	if err != nil {
		return Task{}, &WriteError{Message: fmt.Sprintf("cannot write task: %v", err)}
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return Task{}, &WriteError{Message: fmt.Sprintf("cannot get last inserted id: %v", err)}
	}

	return sb.Get(int(insertedId))
}

func (sb *Sqlite3Buffer) WriteBatch(tasks []Task) ([]Task, error) {
	tx, err := sb.db.Begin()

	if err != nil {
		return nil, &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.WriteBatch]: cannot start transaction: %v", err)}
	}

	stmt, err := tx.Prepare("INSERT INTO tasks (createdAt, updatedAt, name, description, dueDate, isDone, uuid) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return nil, &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.WriteBatch]: cannot prepare statement: %v", err)}
	}

	for _, task := range tasks {
		_, err = stmt.Exec(task.CreatedAt, task.UpdatedAt, task.Name, task.Description, task.DueDate, task.isDone, task.Uuid)

		if err != nil {
			return nil, &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.WriteBatch]: cannot execute statement: %v", err)}
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.WriteBatch]: cannot commit transaction: %v", err)}
	}

	return tasks, nil
}

func (sb *Sqlite3Buffer) Remove(id int) error {
	_, err := sb.db.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err != nil {
		return &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.Remove]: cannot remove task: %v", err)}
	}

	return nil
}

func (sb *Sqlite3Buffer) RemoveBatch(ids []int) error {
	tx, err := sb.db.Begin()

	if err != nil {
		return &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.RemoveBatch]: cannot start transaction: %v", err)}
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE id IN (?)", ids)
	return nil
}

func (sb *Sqlite3Buffer) Update(id int, task Task) (Task, error) {
	_, err := sb.db.Exec(
		"UPDATE tasks SET name = ?, description = ?, dueDate = ?, isDone = ?, uuid = ? WHERE id = ?",
		task.Name, task.Description, task.DueDate, task.isDone, task.Uuid, id,
	)

	if err != nil {
		return Task{}, &WriteError{Message: fmt.Sprintf("[Sqlite3Buffer.Update]: cannot update task: %v", err)}
	}

	return task, nil
}

func (sb *Sqlite3Buffer) Get(id int) (Task, error) {
	row := sb.db.QueryRow("SELECT id, createdAt, updatedAt, name, description, dueDate, isDone, uuid FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.DueDate, &task.isDone, &task.Uuid)
	if err != nil {
		return Task{}, &ReadError{Message: fmt.Sprintf("[Sqlite3Buffer.Get]: cannot get task: %v", err)}
	}

	return task, nil
}

func (sb *Sqlite3Buffer) GetAll() ([]Task, error) {
	rows, err := sb.db.Query("SELECT id, createdAt, updatedAt, name, description, dueDate, isDone, uuid FROM tasks ORDER BY id")

	if err != nil {
		return nil, &ReadError{Message: fmt.Sprintf("[Sqlite3Buffer.GetAll]: cannot get tasks: %v", err)}
	}

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.DueDate, &task.isDone, &task.Uuid)

		if err != nil {
			return nil, &ReadError{Message: fmt.Sprintf("[Sqlite3Buffer.GetAll]: cannot scan row: %v", err)}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (sb *Sqlite3Buffer) GetLatest() (Task, error) {
	row := sb.db.QueryRow("SELECT id, createdAt, updatedAt, name, description, dueDate, isDone, uuid FROM tasks ORDER BY id DESC LIMIT 1")

	var task Task
	err := row.Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.DueDate, &task.isDone, &task.Uuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, &ReadError{Message: fmt.Sprintf("[Sqlite3Buffer.GetLatest]: no tasks found: %v", err)}
		}

		return Task{}, &ReadError{Message: fmt.Sprintf("[Sqlite3Buffer.GetLatest]: cannot get latest task: %v", err)}
	}

	return task, nil
}
