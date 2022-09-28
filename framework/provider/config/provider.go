package config

import (
	"path/filepath"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebConfigProvider struct {
}

func (provider *WebConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewWebConfig
}

func (provider *WebConfigProvider) Boot(c framework.Container) error {
	return nil
}

func (provider *WebConfigProvider) IsDefer() bool {
	return false
}

func (provider *WebConfigProvider) Params(c framework.Container) []any {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []any{c, envFolder, envService.All()}

}

func (provider *WebConfigProvider) Name() string {
	return contract.ConfigKey
}
