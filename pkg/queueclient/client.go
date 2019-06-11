package queueclient

// QueueClient is an interface to deal with a queue of strings
// All QueueClient instances should be able to enqueue, dequeue
// and switch queues
type QueueClient interface {
	Enqueue(name, value string) error
	Dequeue(name string) (string, error)
	SwitchQueue(src, dest string) (string, error)
	String() string
}
