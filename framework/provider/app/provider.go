package app

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

// AppProvider implement ServiceProvider interface
type AppProvider struct {
	BaseFolder string
}

func (provider *AppProvider) Name() string {
	return contract.AppKey
}

func (provider *AppProvider) Params(container framework.Container) []any {
	return []any{container, provider.BaseFolder}
}

func (provider *AppProvider) Register(container framework.Container) framework.NewInstance {
	return NewApp
}

func (provider *AppProvider) Boot(framework.Container) error {
	return nil
}

func (provider *AppProvider) IsDefer() bool {
	return false
}
