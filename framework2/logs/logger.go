package logs

import (
	"io"
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
	for _,opt := range opts {
		opt(l.opt)
	}
}

func (l *logger) Writer() io.Writer {
	return l
}

func (l *logger) Write(data []byte) (int,error) {
	l.entry().write(l.opt.stdLevel, FmtEmptySeparate,*(*string)(unsafe.Pointer(&data)))
	return 0,nil
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

