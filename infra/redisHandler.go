package infra

import (
	"github.com/go-redis/redis"
	"log"
)

// RedisStorager implements the JobHandler interface
type RedisStorager struct {
	*redis.Client
}

// NewRedisHandler creates new redis handler
func NewRedisHandler(addr string) (*RedisStorager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisHandler := new(RedisStorager)
	redisHandler.Client = client

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisHandler, nil
}

// Lookup gets the data based on key
func (s *RedisStorager) Lookup(key string) string {
	val, err := s.Get(key).Result()
	if err != nil {
		log.Println("[RedisStorager.Lookup]: Failed to found key", key, "from redis")
		return ""
	}
	return val
}