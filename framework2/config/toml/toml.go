package toml

import (
	"fmt"
	"github.com/awaketai/goweb/framework2/config"
	"github.com/pelletier/go-toml"
	"os"
)

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
	fmt.Println("1:", ctx)
	return nil, nil
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

func (c configContainer) Set(key, val string) error {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) String(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Strings(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Int(key string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Int64(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Bool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Float(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultString(key string, defaultVal string) string {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultStrings(key string, defaultVal []string) []string {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultInt(key string, defaultVal int) int {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultInt64(key string, defaultVal int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultBool(key string, defaultVal bool) bool {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DefaultFloat(key string, defaultVal float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) DIY(key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) GetSection(section string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) Sub(key string) (config.Configer, error) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) OnChange(key string, fn func(value string)) {
	//TODO implement me
	panic("implement me")
}

func (c configContainer) SaveConfigFile(filename string) error {
	//TODO implement me
	panic("implement me")
}
