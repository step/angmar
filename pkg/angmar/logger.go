package angmar

import (
	"fmt"
	"log"
	"strings"
)

// AngmarLogger is a simple wrapper around a log.Logger
// and provides several convenience methods to log events
// specific to Angmar
type AngmarLogger struct {
	Logger *log.Logger
}

// StartAngmar should be called when Angmar.Start is called.
// It logs the Angmar instance the queue that Angmar listens to.
func (l AngmarLogger) StartAngmar(a Angmar, queueName string) {
	var builder strings.Builder
	builder.WriteString("Starting Angmar...\n")
	builder.WriteString("---\n")
	builder.WriteString(a.String())
	builder.WriteString("Listening to queue: " + queueName + "\n")
	builder.WriteString("---\n")

	l.Logger.Println(builder.String())
}

// ReceivedMessage should be called when a worker picks up an
// Angmar message and before it executes it.
func (l AngmarLogger) ReceivedMessage(workerId int, message AngmarMessage) {
	var builder strings.Builder
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("Received Job...\n")
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	l.Logger.Println(builder.String())
}

// LogError should be called on any Error that occurs within a worker.
// It logs the worker id, the error and the Angmar Message for which
// the error occurred.
func (l AngmarLogger) LogError(workerId int, err error, message AngmarMessage) {
	var builder strings.Builder
	builder.WriteString("Error!\n")
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	builder.WriteString(err.Error())
	l.Logger.Println(builder.String())
}

// TaskPlacedOnQueue should be called when a worker has finished downloading
// and has placed a task on the respective downstream queue. It logs the worker id,
// the Angmar Message and the name of the queue on which the task is placed.
func (l AngmarLogger) TaskPlacedOnQueue(workerId int, message AngmarMessage, qName string) {
	var builder strings.Builder
	builder.WriteString("Task placed on queue\n")
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	builder.WriteString("queue: " + qName + "\n")
	l.Logger.Println(builder.String())
}
