package buffers

type Buffer interface {
	Write(name string) ([]string, error)
	Remove(name int) ([]string, error)
	Get() ([]string, error)
}
