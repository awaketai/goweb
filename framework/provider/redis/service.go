package redis

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/go-redis/redis/v8"
	"sync"
)

type WebRedisService struct {
	container framework.Container
	clients   map[string]*redis.Client
	lock      *sync.RWMutex
}

func NewWebRedisService(params ...any) (any, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}
	return &WebRedisService{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil

}

func (rdService WebRedisService) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	config := GetBaseConfig(rdService.container)
	for _, opt := range option {
		if err := opt(rdService.container, config); err != nil {
			return nil, err
		}
	}

	key := config.UniqKey()
	rdService.lock.RLock()
	if db, ok := rdService.clients[key]; ok {
		rdService.lock.RUnlock()
		return db, nil
	}
	rdService.lock.RUnlock()

	rdService.lock.Lock()
	client := redis.NewClient(config.Options)
	rdService.clients[key] = client
	rdService.lock.Unlock()
	return client, nil
}
