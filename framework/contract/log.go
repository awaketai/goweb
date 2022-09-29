package contract

import (
	"context"
	"io"
	"time"
)

const LogKey = "web:log"

type LogLevel uint32

// panic，表示会导致整个程序出现崩溃的日志信息
// fatal，表示会导致当前这个请求出现提前终止的错误信息
// error，表示出现错误，但是不一定影响后续请求逻辑的错误信息
// warn，表示出现错误，但是一定不影响后续请求逻辑的报警信息
// info，表示正常的日志信息输出debug，表示在调试状态下打印出来的日志信息
// trace，表示最详细的信息，一般信息量比较大，可能包含调用堆栈等信息

const (
	UnknowLevel LogLevel = iota
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Log interface {
	// Panic 表示会导致整个程序出现崩溃的日志信息
	Panic(ctx context.Context, msg string, fields map[string]any)
	// Fatal 表示会导致当前这个请求出现提前终止的错误信息
	Fatal(ctx context.Context, msg string, fields map[string]any)
	// Error 表示出现错误信息，但是不一定影响后续请求逻辑的错误信息
	Error(ctx context.Context, msg string, fields map[string]any)
	// Warn 表示出现错误信息，但是不一定影响后续请求逻辑的报警信息
	Warn(ctx context.Context, msg string, fields map[string]any)
	// Info 正常的日志信息
	Info(ctx context.Context, msg string, fields map[string]any)
	// Debug 调试状态下打印出来的信息
	Debug(ctx context.Context, msg string, fields map[string]any)
	// Trace
	Trace(ctx context.Context, msg string, fields map[string]any)
	// SetLevel
	SetLevel(level LogLevel)

	// SetCtxFielder 从context中获取上下文字段field
	SetCtxFielder(hander CtxFielder)
	// SetFormatter 设置输出格式
	SetFormatter(formatter Formatter)
	// SetOutput 设置输出格式
	SetOutput(out io.Writer)
}

// 日志级别，输出当前日志的级别信息。
// 日志时间，输出当前日志的打印时间。
// 日志简要信息，输出当前日志的简要描述信息，一句话说明日志错误。
// 日志上下文字段，输出当前日志的附带信息。这些字段代表日志打印的上下文。

// CtxFielder 从context获取信息的方法
type CtxFielder func(ctx context.Context) map[string]any

// Formatter 将日志信息组织成字符串的通用方法
type Formatter func(level LogLevel, t time.Time, msg string, fields map[string]any) ([]byte, error)
