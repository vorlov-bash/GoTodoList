package tasks

import (
	"encoding/json"
	"os"
)

type FileBuffer struct {
	FileName string

	// MemBuff is a buffer that holds the data in memory before writing it to the file
	MemBuff MemoryBuffer
}

func writeDataToFile(f *FileBuffer, data []Task) error {
	file, err := os.OpenFile(f.FileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return &WriteError{Message: "os.OpenFile error: " + err.Error()}
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	// json lib automatically marshals the data to json format
	err = enc.Encode(data)

	if err != nil {
		return &WriteError{Message: "json.Encode error: " + err.Error()}
	}

	return nil
}

func getDataFromFile(f *FileBuffer) ([]Task, error) {
	file, err := os.OpenFile(f.FileName, os.O_RDONLY, os.ModePerm)
	defer file.Close()

	if err != nil {
		return nil, &ReadError{Message: "os.Open error: " + err.Error()}
	}

	dec := json.NewDecoder(file)
	var data []Task
	err = dec.Decode(&data)

	if err != nil {
		if err.Error() == "EOF" {
			return []Task{}, nil
		}

		return nil, &ReadError{Message: "json.Decode error: " + err.Error()}
	}

	return data, nil
}

func NewFileBuffer(fileName string) (*FileBuffer, error) {
	// Ensure that the directory exists
	err := os.MkdirAll("tmp", os.ModePerm)

	if err != nil {
		return nil, err
	}

	// Ensure that the file exists
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	fBuff := &FileBuffer{FileName: fileName, MemBuff: MemoryBuffer{data: []Task{}}}

	// Ensure that the buffer is in sync with the file
	data, err := getDataFromFile(fBuff)
	if err != nil {
		return nil, err
	}
	fBuff.MemBuff.data = data
	return fBuff, nil
}

func (fb *FileBuffer) Write(data Task) (Task, error) {
	result, err := fb.MemBuff.Write(data)

	if err != nil {
		return Task{}, err
	}

	err = writeDataToFile(fb, fb.MemBuff.data)
	return result, err
}

func (fb *FileBuffer) WriteBatch(data []Task) ([]Task, error) {
	result, err := fb.MemBuff.WriteBatch(data)

	if err != nil {
		return nil, err
	}

	err = writeDataToFile(fb, result)
	return result, err
}

func (fb *FileBuffer) Remove(id int) error {
	err := fb.MemBuff.Remove(id)

	if err != nil {
		return err
	}

	err = writeDataToFile(fb, fb.MemBuff.data)
	return err
}

func (fb *FileBuffer) RemoveBatch(ids []int) error {
	err := fb.MemBuff.RemoveBatch(ids)

	if err != nil {
		return err
	}

	err = writeDataToFile(fb, fb.MemBuff.data)
	return err
}

func (fb *FileBuffer) Update(id int, data Task) (Task, error) {
	result, err := fb.MemBuff.Update(id, data)

	if err != nil {
		return Task{}, err
	}

	err = writeDataToFile(fb, fb.MemBuff.data)
	return result, err
}

func (fb *FileBuffer) Get(id int) (Task, error) {
	fileData, err := getDataFromFile(fb)

	if err != nil {
		return Task{}, err
	}

	fb.MemBuff.data = fileData
	return fb.MemBuff.Get(id)
}

func (fb *FileBuffer) GetAll() ([]Task, error) {
	fileData, err := getDataFromFile(fb)

	if err != nil {
		return nil, err
	}

	fb.MemBuff.data = fileData
	return fb.MemBuff.GetAll()
}

func (fb *FileBuffer) GetLatest() (Task, error) {
	fileData, err := getDataFromFile(fb)

	if err != nil {
		return Task{}, err
	}

	fb.MemBuff.data = fileData
	return fb.MemBuff.GetLatest()
}
