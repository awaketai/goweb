package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/awaketai/goweb/framework/contract"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

func (conf *WebConfig) loadYAMLConfigFile(folder, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// read file
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		bf = replaceEnvFunc(bf, conf.envMaps)
		c := make(map[string]any)
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		if name == "app" && conf.container.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				service := conf.container.MustMake(contract.AppKey).(contract.App)
				service.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

func (conf *WebConfig) removeYAMLConfig(folder, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		delete(conf.confMaps, name)
		delete(conf.confRaws, name)
	}

	return nil
}
