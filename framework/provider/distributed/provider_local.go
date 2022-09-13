package distributed

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type LocalDistributedProvider struct {
}

func (provider *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

func (provider *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

func (provider *LocalDistributedProvider) IsDefer() bool {
	return false
}

func (provider *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}

func (provider *LocalDistributedProvider) Params(container framework.Container) []any {
	return []any{container}
}
