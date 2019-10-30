package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

type CustomizeRdsStore struct {
	RedisClient *redis.Client
	ExpireAt    time.Duration
}

func GetStore() *CustomizeRdsStore {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	c := new(CustomizeRdsStore)
	c.RedisClient = client
	c.ExpireAt = time.Duration(1 * time.Hour)
	return c
}

func NewStore(expireAt time.Duration) *CustomizeRdsStore {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	c := new(CustomizeRdsStore)
	c.RedisClient = client
	c.ExpireAt = expireAt
	return c
}

func (s CustomizeRdsStore) SetWithOverrideExpire(id string, value string, expireAt time.Duration) {
	err := s.RedisClient.Set(id, value, expireAt).Err()
	if err != nil {
		log.Println(err)
	}
}

func (s CustomizeRdsStore) SetWithoutExpire(id string, value string) {
	err := s.RedisClient.Set(id, value, 0).Err()
	if err != nil {
		log.Println(err)
	}
}

// customizeRdsStore implementing Set method of  Store interface
func (s CustomizeRdsStore) Set(id string, value []byte) {
	err := s.RedisClient.Set(id, string(value), s.ExpireAt).Err()
	if err != nil {
		log.Println(err)
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s CustomizeRdsStore) Get(id string, clear bool) (value []byte) {
	val, err := s.RedisClient.Get(id).Result()
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	if clear {
		err := s.RedisClient.Del(id).Err()
		if err != nil {
			log.Println(err)
			return []byte{}
		}
	}
	return []byte(val)
}
