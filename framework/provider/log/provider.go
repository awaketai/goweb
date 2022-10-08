package log

import (
	"io"
	"strings"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/provider/log/formatter"
	"github.com/awaketai/goweb/framework/provider/log/service"
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

func (provider *WebLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if provider.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return service.NewWebConsoleLog
		}
		cs := tcs.(contract.Config)
		provider.Driver = strings.ToLower(cs.GetString("log.driver"))
	}

	switch provider.Driver {
	case "single":
		return service.NewWebSingleLog
	case "rotate":
		return service.NewWebRotateLog
	case "console":
		return service.NewWebConsoleLog
	case "custom":
		return service.NewWebCustomLog
	default:
		return service.NewWebConsoleLog
	}
}

func (provider *WebLogServiceProvider) Boot(container framework.Container) error {
	return nil
}

func (provider *WebLogServiceProvider) IsDefer() bool {
	return false
}

func (provider *WebLogServiceProvider) Params(container framework.Container) []any {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	if provider.Formatter == nil {
		provider.Formatter = formatter.TextFormatter
		if configService.IsExists("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				provider.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				provider.Formatter = formatter.TextFormatter
			}
		}
	}

	if provider.Level == contract.UnknowLevel {
		provider.Level = contract.InfoLevel
		if configService.IsExists("log.level") {
			provider.Level = logLevel(configService.GetString("log.level"))
		}
	}

	return []any{container, provider.Level, provider.CtxFielder, provider.Formatter, provider.Output}
}

func (provider *WebLogServiceProvider) Name() string {
	return contract.LogKey
}

func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknowLevel
}
