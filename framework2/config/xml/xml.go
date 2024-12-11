package xml

import (
	"errors"
	"github.com/awaketai/goweb/framework2/config"
	"github.com/beego/x2j"
	"os"
	"sync"
)

type Config struct {
}

func (c *Config) Parse(filename string) (config.Configer, error) {
	byts, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return c.ParseData(byts)
}

func (c *Config) ParseData(data []byte) (config.Configer, error) {
	x := &ConfigContainer{
		data: map[string]any{},
	}
	d, err := x2j.DocToMap(string(data))
	if err != nil {
		return nil, err
	}
	v := d["config"]
	if v == nil {
		return nil, errors.New("xml parse should include in <config></config> tags")
	}
	confVal, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("xml parse <config></config> tags should include sub tags")
	}
	x.data = config.ExpandValueEnvForMap(confVal)

	return x, nil
}

type ConfigContainer struct {
	data map[string]any
	sync.Mutex
}

func (c ConfigContainer) Set(key, val string) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) String(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) Strings(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) Int(key string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) Int64(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) Bool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) Float(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultString(key string, defaultVal string) string {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultInt(key string, defaultVal int) int {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	//TODO implement me
	panic("implement me")
}

func (c ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	//TODO implement me
	panic("implement me")
}
