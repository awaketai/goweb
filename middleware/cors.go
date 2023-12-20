package middleware

import (
	"errors"
	"strings"

	"github.com/awaketai/goweb/framework/provider/app"
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
	MaxAge                 int
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
