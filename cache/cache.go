package cache

import (
	"fmt"

	"github.com/go-redis/cache"
	goRedis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// Cache is the fast redis cacher, it serializes & unserializes objects on save/load
var redisCache *cache.Codec

// Init creates the cache instance
func Init(redis *goRedis.Client) {
	redisCache = &cache.Codec{
		Redis: redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	log.Info("Established Redis cache")
}

// Store saves an object in the cache for a certain group using a given key
func Store(group int64, key string, value interface{}) error {
	return redisCache.Set(&cache.Item{
		Key:        getKeyName(group, key),
		Object:     value,
		Expiration: 0,
	})
}

// Fetch returns a value from the cache for a certain group using a given key
func Fetch(group int64, key string, value interface{}) error {
	return redisCache.Get(getKeyName(group, key), &value)
}

func getKeyName(group int64, key string) string {
	return fmt.Sprintf("%v-%v", group, key)
}
