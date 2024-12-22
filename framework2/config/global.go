package config

var globalIns Configer

// from beego
func InitGlobalIns(name string, cfg string) error {
	var err error
	globalIns, err = NewConfig(name, cfg)
	return err
}

func Set(key, val string) error {
	return globalIns.Set(key, val)
}

func String(key string) (string, error) {
	return globalIns.String(key)
}

func Strings(key string) ([]string, error) {
	return globalIns.Strings(key)
}

func Int(key string) (int, error) {
	return globalIns.Int(key)
}

func Int64(key string) (int64, error) {
	return globalIns.Int64(key)
}

func Bool(key string) (bool, error) {
	return globalIns.Bool(key)
}

func Float64(key string) (float64, error) {
	return globalIns.Float(key)
}

// DefaultString support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
func DefaultString(key string, defaultVal string) string {
	return globalIns.DefaultString(key, defaultVal)
}

// DefaultStrings will get string slice
func DefaultStrings(key string, defaultVal []string) []string {
	return globalIns.DefaultStrings(key, defaultVal)
}

func DefaultInt(key string, defaultVal int) int {
	return globalIns.DefaultInt(key, defaultVal)
}

func DefaultInt64(key string, defaultVal int64) int64 {
	return globalIns.DefaultInt64(key, defaultVal)
}

func DefaultBool(key string, defaultVal bool) bool {
	return globalIns.DefaultBool(key, defaultVal)
}

func DefaultFloat(key string, defaultVal float64) float64 {
	return globalIns.DefaultFloat(key, defaultVal)
}

// DIY return the original value
func DIY(key string) (interface{}, error) {
	return globalIns.DIY(key)
}

func GetSection(section string) (map[string]string, error) {
	return globalIns.GetSection(section)
}

func Unmarshaler(prefix string, obj interface{}, opt ...DecodeOption) error {
	return globalIns.Unmarshaler(prefix, obj, opt...)
}

func Sub(key string) (Configer, error) {
	return globalIns.Sub(key)
}

func OnChange(key string, fn func(value string)) {
	globalIns.OnChange(key, fn)
}

func SaveConfigFile(filename string) error {
	return globalIns.SaveConfigFile(filename)
}
