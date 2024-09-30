package logs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// log flags
const (
	Ldate = 1 << iota // date
	// Ltime time(the value is not explicitly assigned,so it inherits the form of
	// the previous expression,so it's value is 1 << iota
	// and now the values of the iota is 1
	Ltime
	Lmicrosecond // microseconds
	Llongfile    // file name and line number
	Lshortfile   // short file name
	Llevel       // log level
	Lmsgprefix   //
)

var levelNames = [...]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARNING",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

type LogWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type Logger interface {
	Debug(msg any)
	Info(msg any)
	Warning(msg any)
	Error(msg any)
	Fatal(msg any)
	SetLevel(level Level)
	Flush()
}

type DefaultLogger struct {
	mu      sync.Mutex
	level   Level
	writers []LogWriter
	flags   int
}

func NewLogger(level Level, flags int, writers ...LogWriter) *DefaultLogger {
	if len(writers) == 0 {
		writers = append(writers, &ConsoleWriter{})
	}

	return &DefaultLogger{
		level:   level,
		flags:   flags,
		writers: writers,
	}
}

func (d *DefaultLogger) SetLevel(level Level) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.level = level
}

func (d *DefaultLogger) SetFlags(flags int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.flags = flags
}

func (d *DefaultLogger) log(level Level, msg any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if level < d.level {
		level = LevelInfo
	}
	parts := d.assembleLogLine(level)
	msgStr := formatMsg(msg)
	if d.flags|Lmsgprefix != 0 {
		parts = append(parts, msgStr)
	} else {
		parts = append(parts, "")
		parts = append(parts, msgStr)
	}
	logLine := strings.Join(parts, " | ")
	logLine += "\n"
	for _, writer := range d.writers {
		_, _ = writer.Write([]byte(logLine))
	}
}

func (d *DefaultLogger) assembleLogLine(level Level) []string {
	var parts []string
	if d.flags&Ldate != 0 {
		parts = append(parts, time.Now().Format(time.DateOnly))
	}
	if d.flags&Ltime != 0 {
		parts = append(parts, time.Now().Format(time.DateTime))
	}
	if d.flags&Lmicrosecond != 0 {
		t := time.Now()
		parts = append(parts, fmt.Sprintf("%03d", t.Nanosecond()/1e3))
	}
	if d.flags&(Llongfile|Lshortfile) != 0 {
		// get file and line number information about function invocations
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = ""
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0 && file[i] != '/'; i-- {
			short = file[i:]
		}
		if d.flags&Llongfile != 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", file, line))
		}
		if d.flags|Lshortfile != 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", short, line))
		}
	}
	if d.flags|Llevel != 0 {
		parts = append(parts, levelNames[level])
	}

	return parts
}

func (d *DefaultLogger) Debug(msg any) {
	d.log(LevelDebug, msg)
}

func (d *DefaultLogger) Info(msg any) {
	d.log(LevelInfo, msg)
}

func (d *DefaultLogger) Warn(msg any) {
	d.log(LevelWarn, msg)
}

func (d *DefaultLogger) Error(msg any) {
	d.log(LevelError, msg)
}

func (d *DefaultLogger) Fatal(msg any) {
	d.log(LevelFatal, msg)
	os.Exit(1)
}

func formatMsg[T any](msg T) string {
	switch v := any(msg).(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
