package buffers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type FileBuffer struct {
	path string
}

func writeDataToFile(f FileBuffer, data []string) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return fmt.Errorf("[writeDataToFile]: os.OpenFile error: %w", err)
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(data)

	return err
}

func NewFileBuffer() (*FileBuffer, error) {
	path := "tmp/tasks.json"

	// Ensure that the directory exists
	err := os.MkdirAll("tmp", os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("[NewFileBuffer]: os.MkdirAll error: %w", err)
	}

	// Ensure that the file exists
	file, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("[NewFileBuffer]: os.OpenFile error: %w", err)
	}
	file.Close()

	return &FileBuffer{path: path}, nil
}

func (f FileBuffer) Write(name string) ([]string, error) {
	data, err := f.Get()
	if err != nil {
		return nil, fmt.Errorf("[FileBuffer].Write(): %w", err)
	}

	data = append(data, name)
	err = writeDataToFile(f, data)

	return data, err
}

func (f FileBuffer) Remove(number int) ([]string, error) {
	data, err := f.Get()

	if err != nil {
		return nil, err
	}

	if number < 1 || number > len(data) {
		return nil, fmt.Errorf("[FileBuffer].Remove(): number out of range")
	}

	index := number - 1
	data = append(data[:index], data[index+1:]...)
	err = writeDataToFile(f, data)

	return data, err
}

func (f FileBuffer) Get() ([]string, error) {
	file, err := os.Open(f.path)

	if err != nil {
		return nil, fmt.Errorf("[FileBuffer].Get(): os.Open file error: %w", err)
	}

	defer file.Close()
	dec := json.NewDecoder(file)
	var data []string
	err = dec.Decode(&data)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) || err.Error() == "EOF" {
			return []string{}, nil
		}
		return nil, fmt.Errorf("[FileBuffer].Get(): json Decode error: %w", err)
	}
	return data, nil
}
