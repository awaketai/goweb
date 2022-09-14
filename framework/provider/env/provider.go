package env

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebEnvProvider struct {
	Folder string
}

func (provider *WebEnvProvider) Register(container framework.Container) framework.NewInstance {
	return NewWebEnv
}

func (provider *WebEnvProvider) IsDefer() bool {
	return false
}

func (provider *WebEnvProvider) Boot(container framework.Container) error {
	app := container.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

func (provider *WebEnvProvider) Params(container framework.Container) []any {
	return []any{provider.Folder}
}

func (provider *WebEnvProvider) Name() string {
	return contract.EnvKey
}
