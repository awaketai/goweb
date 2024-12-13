package xml

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/awaketai/goweb/framework2/config"
	"github.com/beego/x2j"
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

func (c *ConfigContainer) Unmarshaler(prefix string,obj any,opt ...config.DecodeOption) error{
	return nil
}

func (c *ConfigContainer) Sub(key string)(config.Configer,error){
	sub,err := c.sub(key)
	if err != nil { 
		return nil,err
	}

	return &ConfigContainer{
		data: sub,
	},nil

}

func (c *ConfigContainer) sub(key string)(map[string]any,error){
	if key == "" {
		return c.data,nil
	}
	value,ok := c.data[key]
	if !ok {
		return nil,fmt.Errorf("the key is not found:%s",key)
	}
	res,ok := value.(map[string]any)
	if !ok {
		return nil,fmt.Errorf("the value of this key is not a structure:%s",key)
	}

	return res,nil
}

func (c *ConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val

	return nil
}

func (c *ConfigContainer) String(key string) (string, error) {
	if v,ok := c.data[key].(string);ok {
		return v,nil
	}

	return "",nil
}

func (c *ConfigContainer) Strings(key string) ([]string, error) {
	v,err := c.String(key)
	if err != nil || v == "" {
		return nil,err
	}

	return strings.Split(v,";"),nil
}

func (c *ConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key].(string))
}

func (c *ConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key].(string),10,64)
}

func (c *ConfigContainer) Bool(key string) (bool, error) {
	if v := c.data[key];v != nil {
		return config.ParseBool(v)
	}

	return false,fmt.Errorf("not exist key:%q",key)
}

func (c *ConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key].(string),64)
}

func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	v,err := c.String(key)
	if err != nil || v == "" {
		return defaultVal
	}

	return v
}

func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	v,err := c.Strings(key)
	if err != nil || v == nil {
		return defaultVal
	}

	return v
}

func (c *ConfigContainer) DefaultInt(key string, defaultVal int) int {
	v,err := c.Int(key)
	if err != nil {
		return defaultVal
	}

	return v
}

func (c *ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	v,err := c.Int64(key)
	if err != nil {
		return defaultVal
		}

		return v
}

func (c *ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v,err := c.Bool(key)
	if err != nil { 
		return defaultVal
	}

	return v
}

func (c *ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	v,err := c.Float(key )
	if err != nil { 
		return defaultVal
	}

	return v
}
