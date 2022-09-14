package contract

const (
	// EnvProduction 生产环境
	EnvProduction = "production"
	// EnvTesting 测试环境
	EnvTesting = "testing"
	// EnvDevelopment 开发环境
	EnvDevelopment = "development"
	EnvKey = "web:env"
)

type Env interface {
	// AppEnv 获取当前环境
	AppEnv() string
	// IsExists 判断环境变量是否被设置
	IsExists(string) bool
	// Get 获取环境变量
	Get(string) string 
	// All 获取所有环境变量
	All() map[string]string
}