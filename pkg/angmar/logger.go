package angmar

import (
	"fmt"
	"log"
	"strings"
)

type AngmarLogger struct {
	Logger *log.Logger
}

func (l AngmarLogger) StartAngmar(a Angmar, queueName string) {
	var builder strings.Builder
	builder.WriteString("Starting Angmar...\n")
	builder.WriteString("---\n")
	builder.WriteString(a.String())
	builder.WriteString("Listening to queue: " + queueName + "\n")
	builder.WriteString("---\n")

	l.Logger.Println(builder.String())
}

func (l AngmarLogger) ReceivedMessage(workerId int, message AngmarMessage) {
	var builder strings.Builder
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("Received Job...\n")
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	l.Logger.Println(builder.String())
}

func (l AngmarLogger) LogError(workerId int, err error, message AngmarMessage) {
	var builder strings.Builder
	builder.WriteString("Error!\n")
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	builder.WriteString(err.Error())
	l.Logger.Println(builder.String())
}

func (l AngmarLogger) TaskPlacedOnQueue(workerId int, message AngmarMessage, qName string) {
	var builder strings.Builder
	builder.WriteString("Task placed on queue\n")
	workerIdStr := fmt.Sprintf("%d", workerId)
	builder.WriteString("worker id: " + workerIdStr + "\n")
	builder.WriteString(message.String())
	builder.WriteString("queue: " + qName + "\n")
	l.Logger.Println(builder.String())
}
