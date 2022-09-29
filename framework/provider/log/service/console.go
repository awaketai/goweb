package service

import (
	"os"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

// 控制台输出
type WebConsoleLog struct {
	WebLog
}

func NewWebConsoleLog(params ...any) (any, error) {
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	log := &WebConsoleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.c = container
	return log, nil
}
