package storage

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	str, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("REDIS: Error while getting value '%s' by key '%s'", str, key)
	} else {
		log.Printf("REDIS: Got value '%s' by key '%s'", str, key)
	}
	return str, err
}

func (r *RedisStorage) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("REDIS: Error while setting value '%s' by key '%s'", value, key)
	} else {
		log.Printf("REDIS: Set value '%s' by key '%s'", value, key)
	}
	return err
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("REDIS: Error while deleting value by key '%s'", key)
	} else {
		log.Printf("REDIS: Deleted value by key '%s'", key)
	}
	return err
}
