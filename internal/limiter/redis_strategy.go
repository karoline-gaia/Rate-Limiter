package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) Increment(key string, expire time.Duration) (int, error) {
	ctx := context.Background()
	cnt, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if cnt == 1 {
		r.client.Expire(ctx, key, expire)
	}
	return int(cnt), nil
}

func (r *RedisStore) Get(key string) (int, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisStore) SetBlock(key string, duration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key+":block", 1, duration).Err()
}

func (r *RedisStore) IsBlocked(key string) (bool, time.Duration, error) {
	ctx := context.Background()
	res, err := r.client.TTL(ctx, key+":block").Result()
	if err != nil {
		return false, 0, err
	}
	if res > 0 {
		return true, res, nil
	}
	return false, 0, nil
}
