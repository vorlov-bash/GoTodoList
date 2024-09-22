package buffers

type Buffer interface {
	Write(name string) ([]string, error)
	Remove(name string) ([]string, error)
	Get() ([]string, error)
}
