package main

import (
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
	redisConf := redisclient.RedisConf{
		Address:  "localhost:6379",
		Db:       2,
		Password: "",
	}

	redisClient := redisclient.NewDefaultClient(redisConf)

	generator := tarutils.DefaultExtractorGenerator{"/tmp/angmar"}

	file, _ := os.OpenFile("/tmp/angmar.log", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
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

	go a.Start("my_queue", r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
