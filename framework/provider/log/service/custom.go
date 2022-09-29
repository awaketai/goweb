package service

import (
	"io"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

// 自定义输出

type WebCustomLog struct {
	WebLog
}

func NewWebCustomLog(params ...any) (any, error) {
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)
	log := &WebCustomLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(output)
	log.c = container
	return log, nil
}
