package app

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

// WebAppProvider implement ServiceProvider interface
// provider struct name:Web[ProviderName]Provider
// register func name:NewWeb[ProviderName]Service
// contract key name: Web:ProviderName="web:ProviderName"
type WebAppProvider struct {
	BaseFolder string
}

func (provider *WebAppProvider) Name() string {
	return contract.AppKey
}

func (provider *WebAppProvider) Params(container framework.Container) []any {
	return []any{container, provider.BaseFolder}
}

func (provider *WebAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewApp
}

func (provider *WebAppProvider) Boot(framework.Container) error {
	return nil
}

func (provider *WebAppProvider) IsDefer() bool {
	return false
}
