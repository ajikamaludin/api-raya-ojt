package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	rdb *redis.Client
	ctx context.Context
}

var lock = &sync.Mutex{}

func New() *RedisClient {
	return &RedisClient{}
}

func (rd *RedisClient) GetInstance() *redis.Client {
	// fmt.Println("[REDIS] : ", &rd.rdb)
	if rd.rdb == nil {
		configs := configs.GetInstance()
		addr := fmt.Sprintf("%s:%s", configs.Redisconfig.Host, configs.Redisconfig.Port)
		lock.Lock()
		rd.rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: configs.Redisconfig.Password, // no password set
			DB:       0,                            // use default DB
		})
		rd.ctx = context.Background()
		lock.Unlock()
	}
	return rd.rdb
}

func (rd *RedisClient) Get(key string, model interface{}) error {
	redis := rd.GetInstance()
	// fmt.Println("[REDIS][READ] : ", key)
	renotes, err := redis.Get(rd.ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(renotes), &model)
	if err != nil {
		return err
	}

	return nil
}

func (rd *RedisClient) Set(key string, model interface{}, expired time.Duration) error {
	redis := rd.GetInstance()
	// fmt.Println("[REDIS][WRITE] : ", key)
	val, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	err = redis.Set(rd.ctx, key, string(val), expired).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rd *RedisClient) Remove(key string) error {
	redis := rd.GetInstance()
	// fmt.Println("[REDIS][REMOVE] : ", key)
	return redis.Do(rd.ctx, "DEL", key).Err()
}
