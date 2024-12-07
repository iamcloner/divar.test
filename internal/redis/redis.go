package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_SERVER"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       redisDB,
})
var ctx = context.Background()

func Set(key string, value string, exp time.Duration) error {
	if rdb == nil {
		return errors.New("redis Not Initialized")
	}
	err := rdb.Set(ctx, key, value, exp).Err()
	if err != nil {
		return errors.New("failed to set value")
	}
	return nil
}
func Get(key string) (string, error) {
	if rdb == nil {
		return "", errors.New("redis Not Initialized")
	}
	val, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", errors.New("failed to get value")
	}
	return val, nil

}
func Del(key string) error {
	if rdb == nil {
		return errors.New("redis Not Initialized")
	}
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return errors.New("failed to delete value")
	}
	return nil
}
