package main

import (
	"github.com/go-redis/cache"
	goRedis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// Redis is is the redis client struct
var Redis *goRedis.Client

// Cache is the fast redis cacher, it serializes & unserializes objects on save/load
var Cache *cache.Codec

func initRedis() {
	Redis = goRedis.NewClient(&goRedis.Options{
		Addr:     settings.Redis.URI,
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})

	log.Info("Established Redis connection")

	Cache = &cache.Codec{
		Redis: Redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	log.Info("Established Redis cache")
}
