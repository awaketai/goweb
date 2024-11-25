package yaml

import (
	"github.com/awaketai/goweb/framework2/config"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

// from beego config
type Config struct {
}

func (c *Config) Parse(filename string) (*Config, error) {
	//cnf, err := ReadYmlReader(filename)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
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
