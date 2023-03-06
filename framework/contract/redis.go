package contract

import (
	"fmt"
	"github.com/awaketai/goweb/framework"
	"github.com/go-redis/redis/v8"
)

const RedisKey = "web:redis"

type RedisOption func(container framework.Container, config *RedisConfig) error

type Redis interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}

type RedisConfig struct {
	*redis.Options
}

func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
