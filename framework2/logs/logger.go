package logs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"unsafe"
)

const FmtEmptySeparate = ""

type logger struct {
	opt       *options
	mu        sync.Mutex
	entryPool *sync.Pool
}

func NewLogger(opts ...Option) *logger {
	logger := &logger{
		opt: initOptions(opts...),
	}
	logger.entryPool = &sync.Pool{
		New: func() any {
			return entry(logger)
		},
	}

	return logger
}

func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, opt := range opts {
		opt(l.opt)
	}
}

func (l *logger) Writer() io.Writer {
	return l
}

func (l *logger) Write(data []byte) (int, error) {
	l.entry().write(l.opt.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
	return 0, nil
}

func (l *logger) entry() *Entry {
	return l.entryPool.Get().(*Entry)
}

var std = NewLogger()

func StdLogger() *logger {
	return std
}

func Writer() io.Writer {
	return std
}

func (l *logger) Debug(args ...any) {
	l.entry().write(LevelDebug, FmtEmptySeparate, args...)
}

func (l *logger) Info(args ...any) {
	l.entry().write(LevelInfo, FmtEmptySeparate, args...)
}

func (l *logger) Warn(args ...any) {
	l.entry().write(LevelWarn, FmtEmptySeparate, args...)
}

func (l *logger) Error(args ...any) {
	l.entry().write(LevelError, FmtEmptySeparate, args...)
}

func (l *logger) Fatal(args ...any) {
	l.entry().write(LevelFatal, FmtEmptySeparate, args...)
	os.Exit(1)
}

func (l *logger) Panic(args ...any) {
	l.entry().write(LevelPanic, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func (l *logger) Debugf(format string, args ...any) {
	l.entry().write(LevelDebug, format, args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.entry().write(LevelInfo, format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.entry().write(LevelWarn, format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.entry().write(LevelError, format, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.entry().write(LevelFatal, format, args...)
	os.Exit(1)
}

func (l *logger) Panicf(format string, args ...any) {
	l.entry().write(LevelPanic, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// std logger
func Debug(args ...any) {
	std.entry().write(LevelDebug, FmtEmptySeparate, args...)
}

func Info(args ...any) {
	std.entry().write(LevelInfo, FmtEmptySeparate, args...)
}

func Warn(args ...any) {
	std.entry().write(LevelWarn, FmtEmptySeparate, args...)
}

func Error(args ...any) {
	std.entry().write(LevelError, FmtEmptySeparate, args...)
}

func Panic(args ...any) {
	std.entry().write(LevelPanic, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func Fatal(args ...any) {
	std.entry().write(LevelFatal, FmtEmptySeparate, args...)
	os.Exit(1)
}

func Debugf(format string, args ...any) {
	std.entry().write(LevelDebug, format, args...)
}

func Infof(format string, args ...any) {
	std.entry().write(LevelInfo, format, args...)
}

func Warnf(format string, args ...any) {
	std.entry().write(LevelWarn, format, args...)
}

func Errorf(format string, args ...any) {
	std.entry().write(LevelError, format, args...)
}

func Panicf(format string, args ...any) {
	std.entry().write(LevelPanic, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
	std.entry().write(LevelFatal, format, args...)
	os.Exit(1)
}
