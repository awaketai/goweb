package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// core code reference from beego config package
type Configer interface {
	Set(key, val string) error
	String(key string) (string, error)
	Strings(key string) ([]string, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string
	DefaultStrings(key string, defaultVal []string) []string
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
}

type BaseConfiger struct {
	// should support like 'a.b.c'
	reader func(ctx context.Context, key string) (string, error)
}

func NewBaseConfiger(reader func(ctx context.Context, key string) (string, error)) BaseConfiger {
	return BaseConfiger{
		reader: reader,
	}
}

func (b *BaseConfiger) Set(key, val string) error {
	return nil
}

func (b *BaseConfiger) String(key string) (string, error) {
	return b.reader(context.TODO(), key)
}

func (b *BaseConfiger) Strings(key string) ([]string, error) {
	res, err := b.String(key)
	if err != nil || len(res) == 0 {
		return nil, err
	}

	return strings.Split(res, ";"), nil
}

func (b *BaseConfiger) Int(key string) (int, error) {
	res, err := b.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(res)
}

func (b *BaseConfiger) Int64(key string) (int64, error) {
	res, err := b.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(res, 10, 64)
}

func (b *BaseConfiger) Bool(key string) (bool, error) {
	res, err := b.reader(context.TODO(), key)
	if err != nil {
		return false, err
	}

	return ParseBool(res)
}

func (b *BaseConfiger) Float(key string) (float64, error) {
	res, err := b.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(res, 64)
}

func (b *BaseConfiger) DefaultString(key string, defaultVal string) string {
	if res, err := b.String(key); res != "" && err == nil {
		return res
	}

	return defaultVal
}

func (b *BaseConfiger) DefaultStrings(key string, defaultVal []string) []string {
	if res, err := b.Strings(key); len(res) > 0 && err == nil {
		return res
	}

	return defaultVal
}

func (b *BaseConfiger) DefaultInt(key string, defaultVal int) int {
	if res, err := b.Int(key); err == nil {
		return res
	}

	return defaultVal
}

func (b *BaseConfiger) DefaultInt64(key string, defaultVal int64) int64 {
	if res, err := b.Int64(key); err == nil {
		return res
	}

	return defaultVal
}

func (b *BaseConfiger) DefaultBool(key string, defaultVal bool) bool {
	if res, err := b.Bool(key); err == nil {
		return res
	}

	return defaultVal
}

func (b *BaseConfiger) DefaultFloat(key string, defaultVal float64) float64 {
	if re, err := b.Float(key); err == nil {
		return re
	}

	return defaultVal
}

type Config interface {
	Parse(key string) (Configer, error)
	ParseData(data []byte) (Configer, error)
}

func ParseBool(val any) (bool, error) {
	if val == nil {
		return false, fmt.Errorf("parsing <nil>")
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch v {
		case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
			return true, nil
		case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
			return false, nil
		}
	case int8, int32, int64:
		vo := fmt.Sprintf("%d", v)
		switch vo {
		case "1":
			return true, nil
		case "0":
			return false, nil
		}
	case float64:
		switch v {
		case 1.0:
			return true, nil
		case 0.0:
			return false, nil
		}
	}
	return false, fmt.Errorf("parsing %q:invalud value", val)
}

// ExpandValueEnv return value of convert with environment variable
func ExpandValueEnv(val string) (realValue string) {
	realValue = val
	return realValue
}
