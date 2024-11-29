package yaml

import (
	"errors"
	"fmt"
	"github.com/awaketai/goweb/framework2/config"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"sync"
)

func init() {
	config.Register("yaml", &Config{})
}

// from beego config
type Config struct {
}

func (c *Config) Parse(filename string) (config.Configer, error) {
	cnf, err := ReadYmlReader(filename)
	if err != nil {
		return nil, err
	}
	cfg := &ConfigContainer{
		data: cnf,
	}

	return cfg, nil
}

func (c *Config) ParseData(data []byte) (config.Configer, error) {
	cnf, err := parseYAML(data)
	if err != nil {
		return nil, err
	}
	cfg := &ConfigContainer{
		data: cnf,
	}

	return cfg, nil
}

func ReadYmlReader(path string) (map[string]any, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseYAML(buf)
}

func parseYAML(data []byte) (map[string]any, error) {
	cnf := make(map[string]any)
	err := yaml.Unmarshal(data, &cnf)
	if err != nil {
		return nil, err
	}

	cnf = config.ExpandValueEnvForMap(cnf)
	return cnf, nil
}

type ConfigContainer struct {
	data map[string]any
	sync.RWMutex
}

func (c *ConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

func (c *ConfigContainer) String(key string) (string, error) {
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv, nil
		}
	}
	return "", nil
}

func (c *ConfigContainer) Strings(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) Int(key string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) Int64(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) Bool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) Float(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultInt(key string, defaultVal int) int {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (c *ConfigContainer) Unmarshaler(prefix string, obj any, _ ...config.DecodeOption) error {
	sub, err := c.subMap(prefix)
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(sub)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, obj)
}

func (c *ConfigContainer) subMap(key string) (map[string]any, error) {
	tmpData := c.data
	keys := strings.Split(key, ".")
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]any:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return nil, fmt.Errorf("the key is invalid:%s", key)

			}
		}
	}

	return tmpData, nil
}

func (c *ConfigContainer) getData(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := c.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]interface{}:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}
