package logs

import (
	"bytes"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	logger *logger
	Buffer *bytes.Buffer
	Map map[string]any
	Level Level
	Time time.Time
	File string
	Line int
	Func string
	Format string
	Args []any
}

func NewEntry (logger *logger)*Entry {
	return &Entry{
		logger: logger,
		Buffer: &bytes.Buffer{},
		Map: make(map[string]any,5),
	}
}

func (e *Entry) write(level Level,format string,args ...any) {
	if e.logger.opt.level > level {
		return 
	}
	e.Time = time.Now()
	e.Level = level
	e.Format = format
	e.Args = args
	if e.logger.opt.disableCaller {
		if pc,file,line,ok := runtime.Caller(2); !ok {
			e.File = "???"
			e.Func = "???"
		}else {
			e.File = file
			e.Line = line
			e.Func = runtime.FuncForPC(pc).Name()
			e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}
	e.format()
	e.writer()
	e.releas()
}

func (e *Entry) format() {
	_ = e.logger.opt.formatter.Format(e)
}

func (e *Entry) writer() {
	e.logger.mu.Lock()
	_,_ = e.logger.opt.output.Write(e.Buffer.Bytes())
	e.logger.mu.Unlock()
}

func (e *Entry) releas() {
	e.Args = nil
	e.Line = 0
	e.Func = ""
	e.File = ""
	e.Format = ""
	e.Buffer.Reset()
	e.logger.entryPool.Put(e)
}