package env

// from beego config

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var env sync.Map

func init() {
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		env.Store(splits[0], os.Getenv(splits[1]))
	}
	env.Store("GOBIN", GetGOBIN())
	env.Store("GOPATH", GetGOPATH)
}

func Get(key string, defaultValue string) string {
	if val, ok := env.Load(key); ok {
		return val.(string)
	}

	return defaultValue
}

func MustGet(key string) (string, error) {
	if val, ok := env.Load(key); ok {
		return val.(string), nil
	}

	return "", fmt.Errorf("key not found:%v", key)
}

func Set(key string, value string) {
	env.Store(key, value)
}

func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env.Store(key, value)

	return nil
}

func GetAll() map[string]string {
	envs := make(map[string]string, 32)
	env.Range(func(key, value any) bool {
		switch key := key.(type) {
		case string:
			switch val := value.(type) {
			case string:
				envs[key] = val

			}
		}
		return true
	})

	return envs
}

func envFile() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}

	return filepath.Join(dir, "go", "env"), nil
}

func GetRuntimeEnv(key string) (string, error) {
	file, err := envFile()
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("runtime env file is empty")
	}
	var runtimeEnv string
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	envStrings := strings.Split(string(data), "\n")
	for _, envItem := range envStrings {
		envItem = strings.TrimSuffix(envItem, "\r")
		envKeyValue := strings.Split(envItem, "=")
		if len(envKeyValue) == 2 && strings.TrimSpace(envKeyValue[0]) == key {
			runtimeEnv = strings.TrimSpace(envKeyValue[1])
		}
	}

	return runtimeEnv, nil
}

// GetGOBIN return GOBIN environment
func GetGOBIN() string {
	gobin := strings.TrimSpace(Get("GOBIN", ""))
	if gobin == "" {
		var err error
		gobin, err = GetRuntimeEnv("GOBIN")
		if err != nil {
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		if gobin == "" {
			return filepath.Join(build.Default.GOPATH, "bin")
		}

		return gobin
	}

	return gobin
}

// GetGOPATH return GOPATH environment
func GetGOPATH() string {
	gopath := strings.TrimSpace(Get("GOPATH", ""))
	if gopath == "" {
		var err error
		gopath, err = GetRuntimeEnv("GOPATH")
		if err != nil {
			return build.Default.GOPATH
		}
		if gopath == "" {
			return build.Default.GOPATH
		}
		return gopath
	}

	return gopath
}
