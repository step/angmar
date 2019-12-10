package redisclient

import (
	"github.com/step/saurontypes"
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

func (r RedisClient) Add(sName string, entries []saurontypes.Entry) error {
	values := make(map[string]interface{})
	for _, entry := range entries {
		values[entry.Key] = entry.Value
	}
	r.actualClient.XAdd(&redis.XAddArgs{
		Stream: sName,
		ID:     "*",
		Values: values,
	})
	return nil
}

func (r RedisClient) Read(streams []string) []saurontypes.StreamEvent {
	resp := r.actualClient.XRead(&redis.XReadArgs{
		Streams: streams,
		Count:   0,
		Block:   time.Minute,
	})

	if len(resp.Val()) == 0 {
		return []saurontypes.StreamEvent{}
	}
	
	var events []saurontypes.StreamEvent
	for _, val := range resp.Val()[0].Messages {
		streamEvent := saurontypes.StreamEvent{
			ID: val.ID,
			Values: val.Values,
		}
		events = append(events, streamEvent)
	}
	return events
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
