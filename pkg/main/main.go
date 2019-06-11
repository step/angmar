package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/step/angmar/pkg/angmar"

	"github.com/go-redis/redis"
	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/tarutils"
)

type RedisClient struct {
	actualClient *redis.Client
}

func (r RedisClient) Enqueue(name, value string) error {
	r.actualClient.LPush(name, value)
	return nil
}

func (r RedisClient) Dequeue(name string) (string, error) {
	resp := r.actualClient.BRPop(time.Minute, name)
	values, err := resp.Result()
	if err != nil {
		return "", err
	}
	return values[1], err
}

func (r RedisClient) SwitchQueue(src, dest string) (string, error) {
	return "", nil
}

func (r RedisClient) String() string {
	return r.actualClient.String()
}

type DefaultExtractorGenerator struct {
	src string
}

func (d DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	dir := filepath.Join(d.src, args[0], args[1])
	return tarutils.NewDefaultExtractor(dir)
}

func (d DefaultExtractorGenerator) String() string {
	return fmt.Sprintf("DefaultExtractorGenerator: %s\n", d.src)
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	queueClient := RedisClient{client}
	generator := DefaultExtractorGenerator{"/tmp/angmar"}

	file, _ := os.OpenFile("/tmp/angmar.log", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	multiWriter := io.MultiWriter(file, os.Stdout)
	actualLogger := log.New(multiWriter, "--> ", log.LstdFlags)
	logger := angmar.AngmarLogger{Logger: actualLogger}
	a := angmar.Angmar{QueueClient: queueClient,
		Generator:      generator,
		DownloadClient: gh.GithubAPI{Client: http.DefaultClient},
		Logger:         logger}
	r := make(chan bool, 100)
	stop := make(chan bool)
	a.Start("my_queue", r, stop)
}
