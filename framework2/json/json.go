package json

import (
	"encoding/json"
	"fmt"
	"github.com/awaketai/goweb/framework2/config"
	"io"
	"os"
	"sync"
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
	fmt.Println(string(content))

	//return j.ParseData(content), nil
	return nil, nil
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

func (jc *JSONConfigContainer) Set(key, val string) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (jc *JSONConfigContainer) String(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) Strings(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) Int(key string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) Int64(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) Bool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) Float(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultString(key string, defaultVal string) string {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultInt(key string, defaultVal int) int {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (jc *JSONConfigContainer) name() {

}
