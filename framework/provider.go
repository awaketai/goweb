package framework

type NewInstance func(...any) (any, error)

type ServiceProvider interface {
	// Register 在服务容器中注册了一个实例化服务方法
	Register(Container) NewInstance
	// Boot 在调用实例化服务的时候，做一些准备工作
	Boot(Container) error
	// IsDefer 是否在注册的时候实例化服务
	IsDefer() bool
	// Params 传递给NewInstance的参数
	Params(Container) []any
	// Name 服务提供者凭证
	Name() string
}
