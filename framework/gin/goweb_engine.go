package gin

import "github.com/awaketai/goweb/framework"

func (engine *Engine) SetContainer(container framework.Container) {
	engine.container = container

}

func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

func (engine *Engine) GetContainer() framework.Container {
	return engine.container
}
