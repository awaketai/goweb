package contract

import "time"

const ConfigKey = "web:config"

type Config interface {
	// IsExists 检查一个属性是否存在
	IsExists(key string) bool
	Get(key string) any
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetString(key string) string
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	// Load 加载配置到某个对象
	Load(key string, val any) error
}
