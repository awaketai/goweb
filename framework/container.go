package framework

import (
	"fmt"
	"sync"
)

type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已绑定服务提供者
	IsBind(key string) bool
	// Make 根据关键字凭证获取一个服务提供者
	Make(key string) (any, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么panic
	MustMake(key string) any
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例时候非常有用
	MakeNew(key string, params []any) (any, error)
}

type WebContainer struct {
	// 强制要求WebContainer实现Container接口
	Container
	// providers 存储注册的服务提供者
	providers map[string]ServiceProvider
	// instances 存储具体实现
	instances map[string]any
	lock      sync.RWMutex
}

var _ Container = new(WebContainer)

func NewWebContainer() *WebContainer {
	return &WebContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]any),
		lock:      sync.RWMutex{},
	}
}

func (con *WebContainer) Bind(provider ServiceProvider) error {
	key := provider.Name()
	con.lock.Lock()
	con.providers[key] = provider
	con.lock.Unlock()
	if !provider.IsDefer() {
		if err := provider.Boot(con); err != nil {
			return err
		}
		params := provider.Params(con)
		method := provider.Register(con)
		instance, err := method(params...)
		if err != nil {
			return err
		}
		con.lock.Lock()
		con.instances[key] = instance
		con.lock.Unlock()
	}
	return nil
}

func (con *WebContainer) Make(key string) (any, error) {
	return con.make(key, nil, false)

}

func (con *WebContainer) MakeNew(key string, parmas []any) (any, error) {
	return con.make(key, parmas, true)
}

func (con *WebContainer) MustMake(key string) any {
	serv, err := con.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (con *WebContainer) make(key string, params []any, forceNew bool) (any, error) {
	con.lock.RLock()
	defer con.lock.RUnlock()
	provider := con.findServiceProvider(key)
	if provider == nil {
		return nil, fmt.Errorf("contract [" + key + "] not register")

	}
	if forceNew {
		return con.newInstance(provider, params)
	}

	if ins, ok := con.instances[key]; ok {
		return ins, nil
	}
	// 实例化
	inst, err := con.newInstance(provider, params)
	if err != nil {
		return nil, err
	}
	con.instances[key] = inst
	return inst, nil
}

func (con *WebContainer) findServiceProvider(key string) ServiceProvider {
	con.lock.RLock()
	defer con.lock.RUnlock()
	if sp, ok := con.providers[key]; ok {
		return sp
	}
	return nil
}

func (con *WebContainer) newInstance(provider ServiceProvider, params []any) (any, error) {
	if err := provider.Boot(con); err != nil {
		return nil, err
	}

	if len(params) == 0 {
		params = provider.Params(con)
	}
	method := provider.Register(con)
	instance, err := method(params...)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (con *WebContainer) IsBind(key string) bool {
	return con.findServiceProvider(key) != nil
}

func (con *WebContainer) NameList() []string {
	ret := []string{}
	for _, provider := range con.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}
