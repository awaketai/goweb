package contract

const AppKey = "web:app"

type App interface {
	// Version 当前版本
	Version() string
	// BaseFolder 项目基础目录
	BaseFolder() string
	// ConfigFolder config foler
	ConfigFolder() string
	// LogFolder the folder of log
	LogFolder() string
	// ProviderFolder the forlder of provider
	ProviderFolder() string
	// MiddlewareFolder middleware folder
	MiddlewareFolder() string
	// CommandFolder command folder
	CommandFolder() string
	// RuntimeFolder
	RuntimeFolder() string
	TestFolder() string
	// AppFolder 定义业务代码所在目录
	AppFolder() string
	// AppID
	AppID() string
	// LoadAppConfig
	LoadAppConfig(map[string]string)
}
