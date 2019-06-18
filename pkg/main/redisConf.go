package main

import (
	"flag"

	"github.com/step/angmar/pkg/redisclient"
)

var redisAddress string
var redisPassword string
var redisDb int

func init() {
	flag.StringVar(&redisAddress, "redis-address", "localhost:6379", "`address` of Redis host to connect to")
	flag.IntVar(&redisDb, "redis-db", 2, "Redis `database` to transact with")
	flag.StringVar(&redisPassword, "redis-password", "", "`password` for Redis host")
}

func getRedisConf() redisclient.RedisConf {
	return redisclient.RedisConf{
		Address:  redisAddress,
		Password: redisPassword,
		Db:       redisDb,
	}
}
