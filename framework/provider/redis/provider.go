package redis

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebRedisProvider struct {
}

func (rds *WebRedisProvider) Register(container framework.Container) framework.NewInstance {
	return NewWebRedisService
}

func (rds *WebRedisProvider) Boot(container framework.Container) error {
	return nil
}

func (rds *WebRedisProvider) IsDefer() bool {
	return true
}

func (rds *WebRedisProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (rds *WebRedisProvider) Name() string {
	return contract.RedisKey
}
