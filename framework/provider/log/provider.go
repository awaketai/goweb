package log

import (
	"io"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebLogServiceProvider struct {
	framework.ServiceProvider
	Driver     string
	Level      contract.LogLevel
	Formatter  contract.Formatter
	CtxFielder contract.CtxFielder
	// 日志输出信息
	Output io.Writer
}

// func (provider *WebLogServiceProvider) Register(c framework.Container) framework.NewInstance {
// 	if provider.Driver == "" {
// 		tcs, err := c.Make(contract.ConfigKey)
// 		if err != nil {
// 			return
// 		}
// 	}

// }
