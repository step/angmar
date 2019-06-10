package main

import (
	"net/http"
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

type DefaultExtractorGenerator struct {
	src string
}

func (d DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	dir := filepath.Join(d.src, args[0], args[1])
	return tarutils.NewDefaultExtractor(dir)
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	queueClient := RedisClient{client}
	generator := DefaultExtractorGenerator{"/tmp/angmar"}

	a := angmar.Angmar{queueClient, generator, gh.GithubAPI{http.DefaultClient}}
	r := make(chan bool, 4)
	stop := make(chan bool)
	go func() {
		a.Start("my_queue", r, stop)
	}()
	for {
		time.Sleep(time.Microsecond * 100)
	}
}
