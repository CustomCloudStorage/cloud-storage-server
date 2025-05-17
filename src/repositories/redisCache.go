package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CustomCloudStorage/utils"
	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *redisCache {
	return &redisCache{
		client: client,
	}
}

type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Enqueue(ctx context.Context, queue string, value interface{}) error
	Dequeue(ctx context.Context, queue string, dest interface{}) error
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "json marshal for key %s", key)
	}
	if err := r.client.WithContext(ctx).
		Set(key, data, ttl).
		Err(); err != nil {
		return utils.ErrInternal.Wrap(err, "redis SET %s", key)
	}
	return nil
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.WithContext(ctx).Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return utils.ErrNotFound.New("cache miss for key %s", key)
		}
		return utils.ErrInternal.Wrap(err, "redis GET %s", key)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return utils.ErrInternal.Wrap(err, "json unmarshal for key %s", key)
	}
	return nil
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.WithContext(ctx).Exists(key).Result()
	if err != nil {
		return false, utils.ErrInternal.Wrap(err, "redis EXISTS %s", key)
	}
	return n > 0, nil
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.WithContext(ctx).Del(key).Err(); err != nil {
		return utils.ErrInternal.Wrap(err, "redis DEL %s", key)
	}
	return nil
}

func (r *redisCache) Enqueue(ctx context.Context, queue string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "marshal enqueue value")
	}
	if err := r.client.WithContext(ctx).LPush(queue, data).Err(); err != nil {
		return utils.ErrInternal.Wrap(err, fmt.Sprintf("redis LPUSH %s failed", queue))
	}
	return nil
}

func (r *redisCache) Dequeue(ctx context.Context, queue string, dest interface{}) error {
	res, err := r.client.WithContext(ctx).BRPop(0, queue).Result()
	if err != nil {
		return utils.ErrInternal.Wrap(err, fmt.Sprintf("redis BRPOP %s failed", queue))
	}
	var data = []byte(res[1])
	if err := json.Unmarshal(data, dest); err != nil {
		return utils.ErrInternal.Wrap(err, "unmarshal dequeue value")
	}
	return nil
}
