package redisclient

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisConf struct {
	Address  string
	Password string
	Db       int
}

type RedisClient struct {
	actualClient *redis.Client
}

func (r RedisClient) Enqueue(name, value string) error {
	r.actualClient.LPush(name, value)
	return nil
}

func (r RedisClient) Dequeue(name string) (string, error) {
	resp := r.actualClient.BRPop(time.Second*3, name)
	values, err := resp.Result()
	if err != nil {
		return "", err
	}
	return values[1], err
}

func (r RedisClient) Add(sName, key, value string) error {
	entry := make(map[string]interface{})
	entry[key] = value
	r.actualClient.XAdd(&redis.XAddArgs{
		Stream: sName,
		ID: "*",
		Values: entry,
	})
	return nil
}

func (r RedisClient) Read() error  {
	resp := r.actualClient.XRead(&redis.XReadArgs{
		Streams: []string{"eventHub"},
		Count: 0,
		Block: time.Minute,
	})
	
	_, err := resp.Result()

	if err != nil {
		return err
	}

	return nil
}

func (r RedisClient) SwitchQueue(src, dest string) (string, error) {
	return "", nil
}

func (r RedisClient) String() string {
	return r.actualClient.String()
}

func NewDefaultClient(redisConf RedisConf) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Address,
		Password: redisConf.Password, // no password set
		DB:       redisConf.Db,       // use default DB
	})
	return RedisClient{client}
}
