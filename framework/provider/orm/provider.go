package orm

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebGormProvider struct {
}

func (gorm *WebGormProvider) Register(container framework.Container) framework.NewInstance {
	return NewWebOrm
}

func (gorm *WebGormProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer load delay
func (gorm *WebGormProvider) IsDefer() bool {
	return true
}

func (gorm *WebGormProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (gorm *WebGormProvider) Name() string {
	return contract.ORMKey
}
