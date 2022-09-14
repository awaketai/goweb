package env

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/awaketai/goweb/framework/contract"
)

type WebEnv struct {
	// folder .env direcotry
	folder string
	// maps 所有环境变量
	maps map[string]string
}

func NewWebEnv(params ...any) (any, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("new web env params error")
	}

	folder := params[0].(string)
	webEnv := &WebEnv{
		folder: folder,
		// 默认为development环境
		maps: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}
	file := path.Join(folder, ".env")
	fiObj, err := os.Open(file)
	if err == nil {
		defer fiObj.Close()
		// read file
		br := bufio.NewReader(fiObj)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			// parse by equal sign
			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}
			key := string(s[0])
			val := string(s[1])
			webEnv.maps[key] = val
		}
	}
	// 获取当前程序的环境变量，覆盖.env下的变量
	// eg. APP_ENV=testing ./gw env
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		webEnv.maps[pair[0]] = pair[1]
	}
	return webEnv, nil
}

func (en *WebEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

func (en *WebEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

func (en *WebEnv) IsExists(key string) bool {
	_, ok := en.maps[key]
	return ok
}

func (en *WebEnv) All() map[string]string {
	return en.maps
}
