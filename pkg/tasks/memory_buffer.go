package tasks

import (
	"slices"
)

type MemoryBuffer struct {
	// Keep this data SORTED!
	data []Task
}

func NewMemoryBuffer() *MemoryBuffer {
	return &MemoryBuffer{data: []Task{}}
}

func (m *MemoryBuffer) Write(data Task) (Task, error) {
	if len(m.data) > 0 && m.data[len(m.data)-1].Id >= data.Id {
		return Task{}, &WriteError{Message: "can insert only bigger Id"}
	}

	m.data = append(m.data, data)
	return data, nil
}

func (m *MemoryBuffer) WriteBatch(data []Task) ([]Task, error) {
	// Sort input slice by Id
	var sortedData []Task
	for _, task := range data {
		if len(sortedData) == 0 {
			sortedData = append(sortedData, task)
			continue
		}

		// TODO: Try to implement binary algo insert
		for j, _ := range sortedData {
			if task.Id == sortedData[j].Id {
				return nil, &WriteError{Message: "input slice cannot have duplicate id's"}
			}
			sortedData = append(sortedData[:j], append([]Task{task}, sortedData[j:]...)...)
		}
	}

	lastBufferElemId := m.data[len(m.data)-1].Id
	firstSortedElemId := sortedData[len(sortedData)].Id

	// Check if first element of sorted data (the smallest one) is bigger than
	// latest existing element (the greatest one)
	if firstSortedElemId > lastBufferElemId {
		m.data = append(m.data, sortedData...)
		return m.data, nil
	} else if firstSortedElemId == lastBufferElemId {
		return nil, &WriteError{Message: "input slice cannot have duplicate id's"}
	} else {
		return nil, &WriteError{Message: "input slice have greater Id then existing one"}
	}
}

func (m *MemoryBuffer) Remove(id int) error {
	var wasFound bool
	for i, task := range m.data {
		if task.Id == id {
			m.data = append(m.data[:i], m.data[i+1:]...)
			wasFound = true
		}
	}

	if !wasFound {
		return &WriteError{Message: "cannot find task with given id"}
	}

	return nil
}

func (m *MemoryBuffer) RemoveBatch(ids []int) error {
	var batchCopy []Task

	for _, task := range m.data {
		if !slices.Contains(ids, task.Id) {
			batchCopy = append(batchCopy, task)
		}
	}
	m.data = batchCopy

	return nil
}

func (m *MemoryBuffer) Update(id int, data Task) (Task, error) {
	var wasFound bool

	for i, task := range m.data {
		if task.Id == id {
			m.data[i] = data
			wasFound = true
		}
	}

	if !wasFound {
		return Task{}, &WriteError{Message: "cannot find task with given id"}
	}

	return data, nil
}

func (m *MemoryBuffer) Get(id int) (Task, error) {
	for _, task := range m.data {
		if task.Id == id {
			return task, nil
		}
	}
	return Task{}, &ReadError{Message: "cannot find task with given id"}
}

func (m *MemoryBuffer) GetAll() ([]Task, error) {
	return m.data, nil
}

func (m *MemoryBuffer) GetLatest() (Task, error) {
	if len(m.data) == 0 {
		return Task{}, &ReadError{Message: "buffer is empty"}
	}

	return m.data[len(m.data)-1], nil
}
