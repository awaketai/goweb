package service

import (
	"context"
	"io"
	goLog "log"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/provider/log/formatter"
)

type WebLog struct {
	// 日志级别
	level contract.LogLevel
	// 日志格式化方法
	formatter contract.Formatter
	// ctx获取上下文字段
	ctxFielder contract.CtxFielder
	// 输出
	output io.Writer
	// 容器
	c framework.Container
}

func (log *WebLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

func (log *WebLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]any) error {
	// 日志级别
	if !log.IsLevelEnable(level) {
		return nil
	}
	fs := fields
	// 使用ctxFielder获取context中的信息
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		for k, v := range t {
			fs[k] = v
		}
	}
	// 如果绑定了trace服务，获取trace信息
	if log.c.IsBind(contract.TraceKey) {
		tracer := log.c.MustMake(contract.TraceKey).(contract.Trace)
		tc := tracer.GetTrace(ctx)
		if tc != nil {
			maps := tracer.ToMap(tc)
			for k, v := range maps {
				fs[k] = v
			}
		}
	}

	// 将日志信息按照formatter序列化为字符串
	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}
	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}
	if level == contract.PanicLevel {
		goLog.Panicln(string(ct))
		return nil
	}

	log.output.Write(ct)
	log.output.Write([]byte("\r\n"))
	return nil
}

func (log *WebLog) SetOutput(output io.Writer) {
	log.output = output
}

func (log *WebLog) Panic(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

func (log *WebLog) Fatal(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

func (log *WebLog) Error(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (log *WebLog) Warn(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

func (log *WebLog) Info(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

func (log *WebLog) Debug(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

func (log *WebLog) Trace(ctx context.Context, msg string, fields map[string]any) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

func (log *WebLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

func (log *WebLog) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

func (log *WebLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}
