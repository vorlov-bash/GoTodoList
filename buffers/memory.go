package buffers

type MemoryBuffer struct {
	data []string
}

func NewMemoryBuffer() *MemoryBuffer {
	return &MemoryBuffer{
		data: []string{},
	}
}

func (m MemoryBuffer) Write(name string) ([]string, error) {
	m.data = append(m.data, name)
	return m.data, nil
}

func (m MemoryBuffer) Remove(name string) ([]string, error) {
	for i, task := range m.data {
		if task == name {
			m.data = append(m.data[:i], m.data[i+1:]...)
		}
	}

	return m.data, nil
}

func (m MemoryBuffer) Get() ([]string, error) {
	return m.data, nil
}
