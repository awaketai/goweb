package cache

import (
	"context"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/provider/cache/services"
	"strings"
)

type WebCacheProvider struct {
	framework.ServiceProvider
	Driver string
}

func (c *WebCacheProvider) Register(container framework.Container) framework.NewInstance {
	logService := container.MustMake(contract.LogKey).(contract.Log)
	if c.Driver == "" {
		tcs, err := container.Make(contract.ConfigKey)
		if err != nil {
			logService.Error(context.Background(), "WebCacheProvider.Register err:"+err.Error()+"use memory", nil)
			return services.NewMemoryCache
		}
		cs := tcs.(contract.Config)
		c.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	switch c.Driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return services.NewMemoryCache
	default:
		return services.NewMemoryCache

	}
}

func (c *WebCacheProvider) Boot(container framework.Container) error {
	return nil
}

func (c *WebCacheProvider) IsDefer() bool {
	return true
}

func (c *WebCacheProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (c *WebCacheProvider) Name() string {
	return contract.CacheKey
}


