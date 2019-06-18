package main

import (
	"flag"
)

var numberOfWorkers int
var queueName string

func init() {
	flag.IntVar(&numberOfWorkers, "num-of-workers", 5, "`number` of workers that will download in parallel")
	flag.StringVar(&queueName, "queue", "my_queue", "Download `queue` where angmar messages are queued")
}
