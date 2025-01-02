package toml

import (
	"github.com/awaketai/goweb/framework2/config"
	"github.com/pelletier/go-toml"
	"os"
	"strings"
)

// from beego
func init() {
	config.Register("toml", &Config{})
}

// code from beego
const keySeparator = "."

type Config struct {
	tree *toml.Tree
}

func (c *Config) Parse(filename string) (config.Configer, error) {
	ctx, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return c.ParseData(ctx)
}

func (c *Config) ParseData(data []byte) (config.Configer, error) {
	t, err := toml.LoadBytes(data)
	if err != nil {
		return nil, err
	}

	return &configContainer{
		t: t,
	}, nil
}

type configContainer struct {
	t *toml.Tree
}

func (c *configContainer) Set(key, val string) error {
	path := strings.Split(key, keySeparator)
	sub, err := subTree(c.t, path[0:len(path)-1])
	if err != nil {
		return err
	}
	sub.Set(path[len(path)-1], val)
	return nil
}

func (c *configContainer) String(key string) (string, error) {
	res, err := c.get(key)
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", config.KeyNotFoundError
	}

	if str, ok := res.(string); ok {
		return str, nil
	} else {
		return "", config.InvalidValueTypeError
	}
}

func (c *configContainer) Strings(key string) ([]string, error) {
	val, err := c.get(key)
	if err != nil {
		return []string{}, err
	}
	if val == nil {
		return []string{}, config.KeyNotFoundError
	}
	if arr, ok := val.([]interface{}); ok {
		res := make([]string, 0, len(arr))
		for _, ele := range arr {
			if str, ok := ele.(string); ok {
				res = append(res, str)
			} else {
				return []string{}, config.InvalidValueTypeError
			}
		}
		return res, nil
	} else {
		return []string{}, config.InvalidValueTypeError
	}
}

func (c *configContainer) Int(key string) (int, error) {
	val, err := c.Int64(key)
	return int(val), err
}

func (c *configContainer) Int64(key string) (int64, error) {
	res, err := c.get(key)
	if err != nil {
		return 0, err
	}
	if res == nil {
		return 0, config.KeyNotFoundError
	}
	if i, ok := res.(int); ok {
		return int64(i), nil
	} else if i64, ok := res.(int64); ok {
		return i64, nil
	} else {
		return 0, config.InvalidValueTypeError
	}
}

func (c *configContainer) Bool(key string) (bool, error) {
	res, err := c.get(key)
	if err != nil {
		return false, err
	}

	if res == nil {
		return false, config.KeyNotFoundError
	}
	if b, ok := res.(bool); ok {
		return b, nil
	} else {
		return false, config.InvalidValueTypeError
	}
}

func (c *configContainer) Float(key string) (float64, error) {
	res, err := c.get(key)
	if err != nil {
		return 0, err
	}

	if res == nil {
		return 0, config.KeyNotFoundError
	}

	if f, ok := res.(float64); ok {
		return f, nil
	} else {
		return 0, config.InvalidValueTypeError
	}
}

func (c *configContainer) DefaultString(key string, defaultVal string) string {
	res, err := c.get(key)
	if err != nil {
		return defaultVal
	}
	if str, ok := res.(string); ok {
		return str
	} else {
		return defaultVal
	}
}

func (c *configContainer) DefaultStrings(key string, defaultVal []string) []string {
	val, err := c.get(key)
	if err != nil {
		return defaultVal
	}
	if arr, ok := val.([]interface{}); ok {
		res := make([]string, 0, len(arr))
		for _, ele := range arr {
			if str, ok := ele.(string); ok {
				res = append(res, str)
			} else {
				return defaultVal
			}
		}
		return res
	} else {
		return defaultVal
	}
}

func (c *configContainer) DefaultInt(key string, defaultVal int) int {
	return int(c.DefaultInt64(key, int64(defaultVal)))
}

func (c *configContainer) DefaultInt64(key string, defaultVal int64) int64 {
	res, err := c.get(key)
	if err != nil {
		return defaultVal
	}
	if i, ok := res.(int); ok {
		return int64(i)
	} else if i64, ok := res.(int64); ok {
		return i64
	} else {
		return defaultVal
	}
}

func (c *configContainer) DefaultBool(key string, defaultVal bool) bool {
	res, err := c.get(key)
	if err != nil {
		return defaultVal
	}
	if b, ok := res.(bool); ok {
		return b
	} else {
		return defaultVal
	}
}

func (c *configContainer) DefaultFloat(key string, defaultVal float64) float64 {
	res, err := c.get(key)
	if err != nil {
		return defaultVal
	}
	if f, ok := res.(float64); ok {
		return f
	} else {
		return defaultVal
	}
}

func (c *configContainer) DIY(key string) (interface{}, error) {
	return c.get(key)
}

func (c *configContainer) GetSection(section string) (map[string]string, error) {
	val, err := subTree(c.t, strings.Split(section, keySeparator))
	if err != nil {
		return map[string]string{}, err
	}
	m := val.ToMap()
	res := make(map[string]string, len(m))
	for k, v := range m {
		res[k] = config.ToString(v)
	}
	return res, nil
}

func (c *configContainer) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	if len(prefix) > 0 {
		t, err := subTree(c.t, strings.Split(prefix, keySeparator))
		if err != nil {
			return err
		}
		return t.Unmarshal(obj)
	}

	return c.t.Unmarshal(obj)
}

func (c *configContainer) Sub(key string) (config.Configer, error) {
	val, err := subTree(c.t, strings.Split(key, keySeparator))
	if err != nil {
		return nil, err
	}
	return &configContainer{
		t: val,
	}, nil
}

func (c *configContainer) OnChange(key string, fn func(value string)) {

}

func (c *configContainer) SaveConfigFile(filename string) error {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = c.t.WriteTo(f)
	return err
}

func subTree(t *toml.Tree, path []string) (*toml.Tree, error) {
	res := t
	for i := 0; i < len(path); i++ {
		if subTree, ok := res.Get(path[i]).(*toml.Tree); ok {
			res = subTree
		} else {
			return nil, config.InvalidValueTypeError
		}
	}
	if res == nil {
		return nil, config.KeyNotFoundError
	}
	return res, nil
}

func (c *configContainer) get(key string) (interface{}, error) {
	if len(key) == 0 {
		return nil, config.KeyNotFoundError
	}

	segs := strings.Split(key, keySeparator)
	t, err := subTree(c.t, segs[0:len(segs)-1])
	if err != nil {
		return nil, err
	}
	return t.Get(segs[len(segs)-1]), nil
}
