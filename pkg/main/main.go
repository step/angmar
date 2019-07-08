package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/step/angmar/pkg/angmar"
	"github.com/step/angmar/pkg/gh"
)

func main() {
	flag.Parse()

	redisClient := getRedisClient()
	generator := getExtractorGenerator()

	// potentially exits if unable to open file
	file := getLogfile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	logger := getLogger(file)

	ghClient := gh.DefaultGithubAPI()

	a := angmar.NewAngmar(redisClient, generator, ghClient, logger, numberOfWorkers, sourceVolPath)

	r := make(chan bool, 100)
	stop := make(chan bool)

	go a.Start(queueName, r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
