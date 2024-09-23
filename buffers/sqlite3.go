package buffers

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Sqlite3Buffer struct {
	db *sql.DB
}

func NewSqlite3Buffer() (*Sqlite3Buffer, error) {
	path := "tmp/tasks.db"

	err := os.MkdirAll("tmp", os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("[NewSqlite3Buffer]: cannot create directory: %v", err)
	}

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("[NewSqlite3Buffer]: cannot open database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, name TEXT)")
	return &Sqlite3Buffer{db: db}, nil
}

func (s Sqlite3Buffer) Write(name string) ([]string, error) {
	_, err := s.db.Exec("INSERT INTO tasks (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	return s.Get()
}

func (s Sqlite3Buffer) Remove(name int) ([]string, error) {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", name)
	if err != nil {
		return nil, err
	}

	return s.Get()
}

func (s Sqlite3Buffer) Get() ([]string, error) {
	rows, err := s.db.Query("SELECT name FROM tasks")
	if err != nil {
		return nil, err
	}

	var data []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		data = append(data, name)
	}

	return data, nil
}
