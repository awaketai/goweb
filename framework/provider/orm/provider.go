package orm

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type GormProvider struct {
}

func (gorm *GormProvider) Register(container framework.Container) framework.NewInstance {
	return NewWebOrm
}

func (gorm *GormProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer load delay
func (gorm *GormProvider) IsDefer() bool {
	return true
}

func (gorm *GormProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (gorm *GormProvider) Name() string {
	return contract.ORMKey
}
