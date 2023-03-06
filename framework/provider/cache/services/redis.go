package services

import (
	"context"
	"errors"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/provider/redis"
	redisv8 "github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type RedisCache struct {
	container framework.Container
	client    *redisv8.Client
	lock      sync.RWMutex
}

func NewRedisCache(params ...any) (any, error) {
	container := params[0].(framework.Container)
	if !container.IsBind(contract.RedisKey) {
		err := container.Bind(&redis.WebRedisProvider{})
		if err != nil {
			return nil, err
		}
	}

	// 获取redis配置，并实例化redis client
	redisService := container.MustMake(contract.RedisKey).(contract.Redis)
	client, err := redisService.GetClient(redis.WithConfigPath("redis"))
	if err != nil {
		return nil, err
	}
	obj := &RedisCache{
		container: container,
		client:    client,
		lock:      sync.RWMutex{},
	}
	return obj, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	err := r.GetObj(ctx, key, val)
	if err != nil {
		return "", nil
	}
	return val, nil

}

// GetObj 获取某个key对应的对象，对象必须实现https://pkg.go.dev/encoding#BinaryUnMarshaler
func (r *RedisCache) GetObj(ctx context.Context, key string, model any) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv8.Nil) {
		return ErrKeyNotFound
	}
	err := cmd.Scan(model)
	if err != nil {
		return err
	}
	return nil
}

// GetMany 获取某些key对应的值
func (r *RedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redisv8.StringCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	return vals, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

// SetObj 设置某个key和对象到缓存，对象必须实现https://pkg.go.dev/encoding#BinaryMarshaler
func (r *RedisCache) SetObj(ctx context.Context, key string, val any, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *RedisCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipeline := r.client.Pipeline()
	cmds := make([]*redisv8.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipeline.Set(ctx, k, v, timeout))
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func (r *RedisCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DelMany(ctx context.Context, keys []string) error {
	pipeline := r.client.Pipeline()
	cmds := make([]*redisv8.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipeline.Del(ctx, key))
	}
	_, err := pipeline.Exec(ctx)
	return err
}
