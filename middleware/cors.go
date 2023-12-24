package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/awaketai/goweb/framework/gin"
)

// use gin cors

type Config struct {
	AllowAllOrigins bool
	// AllowOrigins a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	AllowOrigins        []string
	AllowOriginFunc     func(origin string) bool
	AllowMethods        []string
	AllowPrivateNetword bool
	AllowHeaders        []string
	// AllowCredentials indicates whether the request can include user credentials
	// like cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// ExposeHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification. Simple methods such as GET and POST are safe to expose to
	ExposeHeaders []string
	// Maxage indicates how long (in seconds) the results of a preflight request
	//can be cached.
	MaxAge                 time.Duration
	AllowWildcard          bool
	AllowBrowserExtensions bool
	// AllowWebSockets allows usage of websocket protocol
	AllowWebSockets bool
	// AllowFiles usage of file:// schema
	AllowFiles               bool
	OptionResponseStatusCode int
}

var (
	DefaultSchemas = []string{
		"http:",
		"https:",
	}
	ExtensionSchemas = []string{
		"chrome-extension://",
		"safari-extension://",
		"moz-extension://",
		"ms-browser-extension://",
	}
	FileSchemas = []string{
		"file://",
	}
	WebSocketSchemas = []string{
		"ws://",
		"wss://",
	}
)

func (c *Config) AddAllowMethods(methods ...string) {
	c.AllowMethods = append(c.AllowMethods, methods...)
}

func (c *Config) AddAllowHeaders(headers ...string) {
	c.AllowHeaders = append(c.AllowHeaders, headers...)
}

func (c *Config) AddExposeHeaders(headers ...string) {
	c.ExposeHeaders = append(c.ExposeHeaders, headers...)
}

func (c *Config) getAllowedSchemas() []string {
	allowedSchmeas := DefaultSchemas
	if c.AllowBrowserExtensions {
		allowedSchmeas = append(allowedSchmeas, "chrome-extension:")
	}
	if c.AllowWebSockets {
		allowedSchmeas = append(allowedSchmeas, WebSocketSchemas...)
	}
	if c.AllowFiles {
		allowedSchmeas = append(allowedSchmeas, FileSchemas...)
	}

	return allowedSchmeas
}

func (c *Config) validateAllowedSchemas(origin string) bool {
	allowedSchemas := c.getAllowedSchemas()
	for _, schema := range allowedSchemas {
		if strings.HasPrefix(origin, schema) {
			return true
		}
	}

	return false
}

func (c *Config) Validate() error {
	if c.AllowAllOrigins && (c.AllowOriginFunc != nil || len(c.AllowOrigins) > 0) {
		return errors.New("conflict settings:all origins are allowed.AllowOriginFunc or AllowOrigins is not needed")
	}
	if !c.AllowAllOrigins && c.AllowOriginFunc == nil && len(c.AllowOrigins) == 0 {
		return errors.New("conflict settings:all origins disabled")
	}
	for _, origin := range c.AllowOrigins {
		if !strings.Contains(origin, "*") && !c.validateAllowedSchemas(origin) {
			return errors.New("bad origin: origins must contain '*' or include " + strings.Join(c.getAllowedSchemas(), ","))
		}
	}

	return nil
}

func (c *Config) parseWildcardRules() [][]string {
	var wRUles [][]string
	if !c.AllowWildcard {
		return wRUles
	}
	for _, o := range c.AllowOrigins {
		if !strings.Contains(o, "*") {
			// find the wildcard
			continue
		}
		if c := strings.Count(o, "*"); c > 1 {
			panic(errors.New("only one * is allowed").Error())
		}
		i := strings.Index(o, "*")
		if i == 0 {
			wRUles = append(wRUles, []string{"*", o[1:]})
			continue
		}
		if i == (len(o) - 1) {
			wRUles = append(wRUles, []string{o[:i-1], "*"})
			continue
		}

		wRUles = append(wRUles, []string{o[:i], o[i+1:]})
	}

	return wRUles
}

func DefaultConfig() Config {
	return Config{
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"HEAD",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// func Default() gin.HandlerFunc {
// 	config := DefaultConfig()
// 	config.AllowAllOrigins = true
// 	return New(config)
// }

// func New(config Config) gin.HandlerFunc {
// 	cors := newCors(config)
// 	return func(c *gin.Context) {
// 		cors.applyCors(c)
// 	}
// }
