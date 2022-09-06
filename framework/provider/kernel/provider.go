package kernel

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/gin"
)

type WebKernelProvider struct {
	HttpEngine *gin.Engine
}

func (provider *WebKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewWebKernelService
}

func (provider *WebKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	provider.HttpEngine.SetContainer(c)
	return nil
}

func (provider *WebKernelProvider) IsDefer() bool {
	return false
}

func (provider *WebKernelProvider) Params(c framework.Container) []any {
	return []any{provider.HttpEngine, "111", 222, "333", "abcd"}
}

func (provider *WebKernelProvider) Name() string {
	return contract.KernelKey
}
