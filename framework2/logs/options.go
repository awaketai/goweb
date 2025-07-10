package logs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var LevelNameMapping = map[Level]string {
	LevelDebug:"DEBUG",
	LevelInfo:"INFO",
	LevelWarn:"WARNING",
	LevelError:"ERROR",
	LevelPanic:"PANIC",
	LevelFatal:"FATAL",
}

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

func (l *Level) unmarshalText(text []byte) bool {
	if len(text) == 0 {
		return false
	}
	levelType := strings.ToLower(string(text))
	switch levelType {
		case "debug":
			*l = LevelDebug
		case "info":
			*l = LevelInfo
		case "warning":
			*l = LevelWarn
		case "error":
			*l = LevelError
		case "panic":
			*l = LevelPanic
		case "fatal":
			*l = LevelFatal
		default:
			return false
	}
	
	return true
}

func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) {
		return fmt.Errorf("unrecognized level:%q", text)
	}
	
	return nil
}

type options struct {
	output io.Writer
	level Level
	stdLevel Level
	formatter Formatter
	disableCaller bool
}

type Option func(*options)

func initOptions(opts ...Option) (*options) {
	o := &options{}
	for _,opt := range opts {
		opt(o)
	}
	
	if o.output == nil {
		o.output = os.Stderr
	}
	if o.formatter == nil {
		
	}
	
	return o
}