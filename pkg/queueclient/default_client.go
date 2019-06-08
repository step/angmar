package queueclient

import (
	"fmt"
	"sync"
)

// InMemoryQueue is an in memory implementation of a queue
// It is "thread safe" as the Enqueue and Dequeue operations
// are made mutually exclusive with a mutex.
type inMemoryQueue struct {
	queues map[string][]string
	mux    sync.Mutex
}

// Enqueue places the given value on the queue specified
// by name. No error is returned since the value is placed
// in memory
func (q inMemoryQueue) Enqueue(name, value string) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queues[name] = append(q.queues[name], value)
	return nil
}

// Dequeue pops a value from the queue specified by name
// An error is returned if Dequeue is attempted on an
// empty queue.
func (q inMemoryQueue) Dequeue(name string) (string, error) {
	queue := q.queues[name]
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(queue) == 0 {
		return "", fmt.Errorf("queue %s length is 0", name)
	}
	value := queue[0]
	q.queues[name] = queue[1:]
	return value, nil
}

// SwitchQueue pops a value from src and pushes it to destination
// and returns the value and/or a potential error. It is similar
// to BRPOPLPUSH
func (q inMemoryQueue) SwitchQueue(src, dest string) (string, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	queue := q.queues[src]
	if len(queue) == 0 {
		return "", fmt.Errorf("queue %s length is 0", src)
	}
	value := queue[0]
	q.queues[src] = queue[1:]
	q.queues[dest] = append(q.queues[dest], value)
	return value, nil
}

// NewDefaultClient provides an in memory implementation of a queue
func NewDefaultClient() QueueClient {
	return inMemoryQueue{queues: make(map[string][]string)}
}
