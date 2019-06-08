package queueclient

type QueueClient interface {
	Enqueue(name, value string) error
	Dequeue(name string) (string, error)
}
