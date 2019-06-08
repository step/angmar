package main

import (
	"net/http"
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
	return values[1], err
}

type DefaultExtractorGenerator struct {
	src string
}

func (d DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	return tarutils.NewDefaultExtractor(d.src)
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
	r := make(chan bool)
	stop := make(chan bool)
	go func() {
		a.Start("my_queue", r, stop)
	}()
	<-r
	<-r
	stop <- true
	for {
	}
}
