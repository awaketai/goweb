package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/awaketai/goweb/framework2/config"
	"github.com/mitchellh/mapstructure"
)

// from beego
type JSONConfig struct {
}

func (j *JSONConfig) Parse(filename string) (config.Configer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return j.ParseData(content)
}

func (j *JSONConfig) ParseData(data []byte) (config.Configer, error) {
	x := &JSONConfigContainer{
		data: map[string]any{},
	}
	err := json.Unmarshal(data, &x.data)
	if err != nil {
		var wrappingArr []any
		innerErr := json.Unmarshal(data, &wrappingArr)
		if innerErr != nil {
			return nil, err
		}
		x.data["rootArray"] = wrappingArr
	}

	return x, nil
}

type JSONConfigContainer struct {
	data map[string]any
	sync.RWMutex
}

func (jc *JSONConfigContainer) sub(key string) (map[string]any, error) {
	if key == "" {
		return jc.data, nil
	}
	value, ok := jc.data[key]
	if !ok {
		return nil, fmt.Errorf("key:%s not exist", key)
	}
	res, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("value:%s not exist", key)
	}

	return res, nil
}

func (jc *JSONConfigContainer) getData(key string) any {
	if key == "" {
		return nil
	}
	jc.RLock()
	defer jc.RUnlock()
	sectionKeys := strings.Split(key, "::")
	if len(sectionKeys) >= 2 {
		curVal, ok := jc.data[sectionKeys[0]]
		if !ok {
			return nil
		}
		for _, secKey := range sectionKeys[1:] {
			if vo, ok := curVal.(map[string]any); ok {
				if curVal, ok = vo[secKey]; !ok {
					return nil
				}
			}
		}

		return curVal
	}
	if v, ok := jc.data[key]; ok {
		return v
	}

	return nil
}

func (jc *JSONConfigContainer) Set(key, val string) error {
	jc.Lock()
	defer jc.Unlock()
	jc.data[key] = val
	return nil
}

func (jc *JSONConfigContainer) String(key string) (string, error) {
	val := jc.getData(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v, nil
		}
	}

	return "", nil
}

func (jc *JSONConfigContainer) Strings(key string) ([]string, error) {
	stringVal, err := jc.String(key)
	if err != nil || stringVal == "" {
		return nil, err
	}

	return strings.Split(stringVal, ";"), nil
}

func (jc *JSONConfigContainer) Int(key string) (int, error) {
	val := jc.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		} else if v, ok := val.(string); ok {
			return strconv.Atoi(v)
		}

		return 0, errors.New("not valid value")
	}

	return 0, errors.New("not valid value")
}

func (jc *JSONConfigContainer) Int64(key string) (int64, error) {
	val := jc.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		}
		return 0, errors.New("not int64 value")
	}
	return 0, errors.New("not exist key:" + key)
}

func (jc *JSONConfigContainer) Bool(key string) (bool, error) {
	val := jc.getData(key)
	if val != nil {
		return config.ParseBool(val)
	}
	return false, fmt.Errorf("not exist key: %q", key)
}

func (jc *JSONConfigContainer) Float(key string) (float64, error) {
	val := jc.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		}
		return 0.0, errors.New("not float64 value")
	}
	return 0.0, errors.New("not exist key:" + key)
}

func (jc *JSONConfigContainer) DefaultString(key string, defaultVal string) string {
	if v, err := jc.String(key); v != "" && err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	if v, err := jc.Strings(key); v != nil && err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DefaultInt(key string, defaultVal int) int {
	if v, err := jc.Int(key); err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := jc.Int64(key); err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	if v, err := jc.Bool(key); err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := jc.Float(key); err == nil {
		return v
	}

	return defaultVal
}

func (jc *JSONConfigContainer) DIY(key string) (any, error) {
	if v, ok := jc.data[key]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("not exist key:%q", key)
}

func (jc *JSONConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := jc.data[section]; ok {
		return v.(map[string]string), nil
	}
	return nil, errors.New("nonexist section " + section)
}

func (jc *JSONConfigContainer) OnChange(key string, fn func(value string)) {

}

func (jc *JSONConfigContainer) SaveConfigFile(filename string) error {
	return nil
}

func (jc *JSONConfigContainer) Sub(key string) (config.Configer, error) {
	sub, err := jc.sub(key)
	if err != nil {
		return nil, err
	}

	return &JSONConfigContainer{
		data: sub,
	}, nil
}

func (jc *JSONConfigContainer) Unmarshaler(prefix string, obj any, opt ...config.DecodeOption) error {
	sub, err := jc.sub(prefix)
	if err != nil {
		return err
	}
	return mapstructure.Decode(sub, obj)
}
