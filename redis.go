package main

import (
	goRedis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// Redis is is the redis client struct
var Redis *goRedis.Client

func initRedis() {
	Redis = goRedis.NewClient(&goRedis.Options{
		Addr:     settings.Redis.URI,
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})

	log.Info("Established Redis connection")
}
