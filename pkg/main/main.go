package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/step/angmar/pkg/angmar"
	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/redisclient"
	"github.com/step/angmar/pkg/tarutils"
)

func main() {
	if os.Args[1] == "help" {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	redisConf := getRedisConf()
	redisClient := redisclient.NewDefaultClient(redisConf)

	generator := tarutils.DefaultExtractorGenerator{Src: sourceVolPath}

	file, err := os.OpenFile(getLogfileName(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	multiWriter := io.MultiWriter(file, os.Stdout)

	actualLogger := log.New(multiWriter, "--> ", log.LstdFlags)
	logger := angmar.AngmarLogger{Logger: actualLogger}

	a := angmar.Angmar{
		QueueClient:    redisClient,
		Generator:      generator,
		DownloadClient: gh.GithubAPI{Client: http.DefaultClient},
		Logger:         logger,
	}

	r := make(chan bool, 100)
	stop := make(chan bool)

	go a.Start(queueName, r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
