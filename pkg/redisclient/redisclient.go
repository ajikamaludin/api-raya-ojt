package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ajikamaludin/go-fiber-rest/app/configs"
	"github.com/go-redis/redis/v8"
)

var lock = &sync.Mutex{}
var rdb *redis.Client
var ctx = context.Background()

func GetInstance() *redis.Client {
	// fmt.Println("[REDIS] : ", &rdb)
	if rdb == nil {
		configs := configs.GetInstance()
		addr := fmt.Sprintf("%s:%s", configs.Redisconfig.Host, configs.Redisconfig.Port)
		lock.Lock()
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: configs.Redisconfig.Password, // no password set
			DB:       0,                            // use default DB
		})
		lock.Unlock()
	}
	return rdb
}

func Get(key string, model interface{}) error {
	redis := GetInstance()
	// fmt.Println("[REDIS][READ] : ", key)
	renotes, err := redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(renotes), &model)
	if err != nil {
		return err
	}

	return nil
}

func Set(key string, model interface{}, expired time.Duration) error {
	redis := GetInstance()
	// fmt.Println("[REDIS][WRITE] : ", key)
	val, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	err = redis.Set(ctx, key, string(val), expired).Err()
	if err != nil {
		return err
	}

	return nil
}

func Remove(key string) error {
	redis := GetInstance()
	// fmt.Println("[REDIS][REMOVE] : ", key)
	return redis.Do(ctx, "DEL", key).Err()
}
