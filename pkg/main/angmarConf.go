package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/step/angmar/pkg/angmar"
)

var numberOfWorkers int
var queueName string

func init() {
	flag.IntVar(&numberOfWorkers, "num-of-workers", 5, "`number` of workers that will download in parallel")
	flag.StringVar(&queueName, "queue", "my_queue", "Download `queue` where angmar messages are queued")
}

func getLogger(file *os.File) angmar.AngmarLogger {
	multiWriter := io.MultiWriter(file, os.Stdout)

	actualLogger := log.New(multiWriter, "--> ", log.LstdFlags)
	return angmar.AngmarLogger{Logger: actualLogger}
}
