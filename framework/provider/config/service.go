package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

type WebConfig struct {
	container framework.Container
	folder    string
	keyDelim  string // 路径分割符，默认为点
	lock      sync.RWMutex
	envMaps   map[string]string // 所有的环境变量
	confMaps  map[string]any    // 配置文件结构，key:文件名
	confRaws  map[string][]byte // 配置文件原始信息
}

var _ contract.Config = new(WebConfig)

func NewWebConfig(params ...any) (any, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("the params index zero must be the framework.Container instance,index first must be the env config folder and the second must be the env config")
	}

	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("the env folder [" + envFolder + "] not exists:" + err.Error())
	}

	webConf := &WebConfig{
		container: container,
		folder:    envFolder,
		envMaps:   envMaps,
		confMaps:  make(map[string]any, 0),
		confRaws:  make(map[string][]byte, 0),
		keyDelim:  ".",
		lock:      sync.RWMutex{},
	}

	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, err
	}
	// 后续优化为interface
	for _, file := range files {
		fileName := file.Name()
		err := webConf.loadYAMLConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	// 监控配置文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}

		}()
		for {
			select {
			case ev := <-watch.Events:
				{
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("create file:", ev.Name)
						webConf.loadYAMLConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("write file:", ev.Name)
						webConf.loadYAMLConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("remove file:", ev.Name)
						webConf.loadYAMLConfigFile(folder, fileName)
					}

				}
			case err := <-watch.Errors:
				{
					log.Println("fsnotify watch err:", err)
					return
				}

			}

		}
	}()

	return webConf, nil
}

// replaceEnvFunc 替换配置文件中的env()配置
func replaceEnvFunc(content []byte, maps map[string]string) []byte {
	if len(maps) == 0 {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

func searchConfig(source map[string]any, path []string) any {
	if len(path) == 0 {
		return source
	}
	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}
		switch next.(type) {
		case map[any]any:
			return searchConfig(cast.ToStringMap(next), path[1:])
		case map[string]any:
			return searchConfig(next.(map[string]any), path[1:])
		default:
			return nil
		}
	}
	return nil
}

func (conf *WebConfig) find(key string) any {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchConfig(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func (conf *WebConfig) IsExists(key string) bool {
	return conf.find(key) != nil
}

func (conf *WebConfig) Get(key string) any {
	return conf.find(key)
}

func (conf *WebConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *WebConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

func (conf *WebConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

func (conf *WebConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

func (conf *WebConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

func (conf *WebConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

func (conf *WebConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

func (conf *WebConfig) GetStringMap(key string) map[string]any {
	return cast.ToStringMap(conf.find(key))
}

func (conf *WebConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

func (conf *WebConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *WebConfig) Load(key string, val any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(conf.find(key))
}
