package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect(addr string) error {
	Client = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Println("connected to redis")
	return nil
}

func Close() {
	if Client != nil {
		Client.Close()
	}
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = Client.Set(ctx, key, data, ttl).Err()
	if err == nil {
		log.Printf("[CACHE SET] %s", key)
	}
	return err
}

func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("[CACHE MISS] %s", key)
		return err
	}
	log.Printf("[CACHE HIT] %s", key)
	return json.Unmarshal(data, dest)
}

func Delete(ctx context.Context, key string) error {
	return Client.Del(ctx, key).Err()
}
