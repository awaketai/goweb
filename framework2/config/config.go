package config

// from beego config
import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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

	// DIY return the original value
	DIY(key string) (any, error)

	GetSection(section string) (map[string]string, error)

	Unmarshaler(prefix string, obj any, opt ...DecodeOption) error
	Sub(key string) (Configer, error)
	OnChange(key string, fn func(value string))
	SaveConfigFile(filename string) error
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

var adapters = make(map[string]Config)

func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config:Register adapter is nil")
	}
	if _, ok := adapters[strings.ToLower(name)]; ok {
		panic("config:Register adapter already registered " + name)
	}
	adapters[strings.ToLower(name)] = adapter
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
		v = strings.ToLower(v)
		switch v {
		case "1", "t", "true", "yes", "y", "on":
			return true, nil
		case "0", "f", "false", "no", "n", "off":
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
func ExpandValueEnv(value string) (realValue string) {
	realValue = value
	vLen := len(value)
	// 3 = ${}
	if vLen < 3 {
		return
	}
	// Need start with "${" and end with "}", then return.
	if value[0] != '$' || value[1] != '{' || value[vLen-1] != '}' {
		return
	}

	key := ""
	defaultV := ""
	// value start with "${"
	for i := 2; i < vLen; i++ {
		if value[i] == '|' && (i+1 < vLen && value[i+1] == '|') {
			key = value[2:i]
			defaultV = value[i+2 : vLen-1] // other string is default value.
			break
		} else if value[i] == '}' {
			key = value[2:i]
			break
		}
	}

	realValue = os.Getenv(key)
	if realValue == "" {
		realValue = defaultV
	}

	return
}

// ExpandValueEnvForMap convert all string value with environment variable.
func ExpandValueEnvForMap(m map[string]any) map[string]any {
	for k, v := range m {
		switch value := v.(type) {
		case string:
			m[k] = ExpandValueEnv(value)
		case map[string]any:
			m[k] = ExpandValueEnvForMap(value)
		case map[string]string:
			for k2, v2 := range value {
				value[k2] = ExpandValueEnv(v2)
			}
			m[k] = value
		case map[any]any:
		tmp := make(map[string]any, len(value))
			for k2, v2 := range value {
				tmp[k2.(string)] = v2
			}
			m[k] = ExpandValueEnvForMap(tmp)
		}
	}
	return m
}

func NewConfig(adapterName, filename string) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("unknown adapter %q", adapter)
	}

	return adapter.Parse(filename)
}

func NewConfigData(adapterName string, data []byte) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("unknown adapter %q", adapter)
	}

	return adapter.ParseData(data)
}

// ToString converts values of any type to string.
func ToString(x any) string {
	switch y := x.(type) {

	// Handle dates with special logic
	// This needs to come above the fmt.Stringer
	// test since time.Time's have a .String()
	// method
	case time.Time:
		return y.Format("A Monday")

	// Handle type string
	case string:
		return y

	// Handle type with .String() method
	case fmt.Stringer:
		return y.String()

	// Handle type with .Error() method
	case error:
		return y.Error()

	}

	// Handle named string type
	if v := reflect.ValueOf(x); v.Kind() == reflect.String {
		return v.String()
	}

	// Fallback to fmt package for anything else like numeric types
	return fmt.Sprint(x)
}

type DecodeOption func(options decodeOptions)
type decodeOptions struct{}
