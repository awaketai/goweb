package app

import (
	"fmt"
	"path/filepath"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/google/uuid"
)

type App struct {
	container  framework.Container
	baseFolder string
	appID      string // 当前app唯一id

	configMap map[string]string // 配置加载
}

var _ contract.App = new(App)

func (app App) Version() string {
	return "0.0.3"
}

func NewApp(params ...any) (any, error) {
	if len(params) != 2 {
		return App{}, fmt.Errorf("param error,the index zero will be container and index one will be baseFolder")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return App{
		baseFolder: baseFolder,
		container:  container,
		appID:      uuid.NewString(),
		configMap:  make(map[string]string),
	}, nil
}

func (app App) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}
	// 获取当前路径
	currentDir, err := util.GetExecDir()
	if err != nil {
		panic("get framework execute dir err:" + err.Error())
	}
	return currentDir
}

// StorageFolder app/storage
func (app App) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

// LogFolder app/storage/log
func (app App) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

// HttpFolder app/http
func (app App) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "http")
}

// ConsoleFolder app/consule
func (app App) ConsoleFolder() string {
	if val, ok := app.configMap["confole_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "console")
}

// ConfigFolder app/config
func (app App) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

// ProviderFolder app/provider
func (app App) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "provider")
}

// MiddlewareFolder app/http/middleware
func (app App) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

// CommandFolder app/console/command
func (app App) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

// RuntimeFolder app/storage/runtime
func (app App) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

// TestFolder app/test
func (app App) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

func (app App) AppFolder() string {
	if val, ok := app.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}

func (app App) AppID() string {
	return app.appID
}

func (app App) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}
