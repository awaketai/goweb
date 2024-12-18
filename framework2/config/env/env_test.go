package env_test

import (
	"github.com/awaketai/goweb/framework2/config/env"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetEnv(t *testing.T) {
	key := "TEST_KEY"
	defaultValue := "default"
	value := "test_value"

	// 清除环境变量，确保测试独立性
	os.Unsetenv(key)

	// 测试 Get 方法的默认值逻辑
	result := env.Get(key, defaultValue)
	assert.Equal(t, defaultValue, result, "Default value should be returned")

	// 设置环境变量并验证
	env.Set(key, value)
	result = env.Get(key, defaultValue)
	assert.Equal(t, value, result, "Set value should be returned")
}

func TestMustGet(t *testing.T) {
	key := "MUST_GET_KEY"
	value := "must_get_value"

	// 清除环境变量，确保测试独立性
	os.Unsetenv(key)

	// 测试未设置时返回错误
	_, err := env.MustGet(key)
	assert.Error(t, err, "MustGet should return an error if the key is not found")

	// 设置环境变量并验证
	env.Set(key, value)
	result, err := env.MustGet(key)
	assert.NoError(t, err, "MustGet should not return an error for existing keys")
	assert.Equal(t, value, result, "MustGet should return the correct value")
}

func TestMustSet(t *testing.T) {
	key := "MUST_SET_KEY"
	value := "must_set_value"

	// 测试设置环境变量
	err := env.MustSet(key, value)
	assert.NoError(t, err, "MustSet should not return an error")
	result := os.Getenv(key)
	assert.Equal(t, value, result, "Environment variable should be set correctly")
}

func TestGetAll(t *testing.T) {
	key := "GET_ALL_KEY"
	value := "get_all_value"

	// 设置一个测试环境变量
	env.Set(key, value)

	// 验证 GetAll 方法
	allEnvs := env.GetAll()
	assert.Contains(t, allEnvs, key, "GetAll should include the test key")
	assert.Equal(t, value, allEnvs[key], "GetAll should return the correct value for the test key")
}

func TestGetRuntimeEnv(t *testing.T) {
	// 创建临时目录和文件模拟运行时环境文件
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "env")
	err := os.WriteFile(envFile, []byte("RUNTIME_KEY=runtime_value\n"), 0644)
	assert.NoError(t, err, "Runtime environment file should be created successfully")

	// 设置 GOENV 环境变量指向临时文件
	err = os.Setenv("GOENV", envFile)
	assert.NoError(t, err, "GOENV should be set without error")

	// 测试读取运行时环境变量
	result, err := env.GetRuntimeEnv("RUNTIME_KEY")
	assert.NoError(t, err, "GetRuntimeEnv should not return an error")
	assert.Equal(t, "runtime_value", result, "GetRuntimeEnv should return the correct value")
}

func TestGetGOBIN(t *testing.T) {
	// 设置和测试 GOBIN 变量
	os.Setenv("GOBIN", "/test/gobin")
	assert.Equal(t, "/test/gobin", env.GetGOBIN(), "GOBIN should return the correct value")
}

func TestGetGOPATH(t *testing.T) {
	// 设置和测试 GOPATH 变量
	os.Setenv("GOPATH", "/test/gopath")
	assert.Equal(t, "/test/gopath", env.GetGOPATH(), "GOPATH should return the correct value")
}
