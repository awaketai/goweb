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
	v, err := c.String(key)
	if v == "" || err != nil {
		return nil, err
	}
	return strings.Split(v, ";"), nil
}

func (c *ConfigContainer) Int(key string) (int, error) {
	if v, err := c.getData(key); err == nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}

	return 0, errors.New("not int value")
}

func (c *ConfigContainer) Int64(key string) (int64, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	default:
		return 0, errors.New("not int or int64 value")
	}
}

func (c *ConfigContainer) Bool(key string) (bool, error) {
	v, err := c.getData(key)
	if err != nil {
		return false, err
	}
	return config.ParseBool(v)
}

func (c *ConfigContainer) Float(key string) (float64, error) {
	if v, err := c.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
}

func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	v, err := c.String(key)
	if v == "" || err != nil {
		return defaultVal
	}
	return v
}

func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	v, err := c.Strings(key)
	if v == nil || err != nil {
		return defaultVal
	}
	return v
}

func (c *ConfigContainer) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultVal
	}
	return v
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

func (c *ConfigContainer) Sub(key string) (config.Configer, error) {
	sub, err := c.subMap(key)
	if err != nil {
		return nil, err
	}

	return &ConfigContainer{
		data: sub,
	}, nil
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

func (c *ConfigContainer) OnChange(_ string, _ func(value string)) {
	// do nothing
}

func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		switch val := v.(type) {
		case map[string]any:
			res := make(map[string]string, len(val))
			for k2, v2 := range val {
				res[k2] = fmt.Sprintf("%v", v2)
			}
			return res, nil
		case map[string]string:
			return val, nil
		default:
			return nil, fmt.Errorf("the section is invalid:%s", section)
		}
	}

	return nil, errors.New("section not found")
}

func (c *ConfigContainer) SaveConfigFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	return err
}

func (c *ConfigContainer) DIY(key string) (any, error) {
	return c.getData(key)
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
