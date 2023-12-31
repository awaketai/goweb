package middleware

import (
	"errors"
	"net/http"
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

type cors struct {
	allowAllOrigins          bool
	allowCredentials         bool
	allowOriginFunc          func(string) bool
	allowOrigins             []string
	normalHeaders            http.Header
	preflightHeaders         http.Header
	wildcardOrigins          [][]string
	optionResponseStatusCode int
}

func newCors(config Config) *cors {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}
	for _, origin := range config.AllowOrigins {
		if origin == "*" {
			config.AllowAllOrigins = true
		}
	}
	if config.OptionResponseStatusCode == 0 {
		config.OptionResponseStatusCode = http.StatusNoContent
	}

	return &cors{
		allowOriginFunc:          config.AllowOriginFunc,
		allowAllOrigins:          config.AllowAllOrigins,
		allowCredentials:         config.AllowCredentials,
		allowOrigins:             normalize(config.AllowOrigins),
		normalHeaders:            generateNormalHeaders(config),
		preflightHeaders:         generatePreflightHeaders(config),
		wildcardOrigins:          config.parseWildcardRules(),
		optionResponseStatusCode: config.OptionResponseStatusCode,
	}
}

func (cors *cors) applyCors(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return
	}
	host := c.Request.Host

	if origin == "http://"+host || origin == "https://"+host {
		// request is not a CORS request but have origin header.
		// for example, use fetch api
		return
	}

	if !cors.validateOrigin(origin) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if c.Request.Method == "OPTIONS" {
		cors.handlePreflight(c)
		defer c.AbortWithStatus(cors.optionResponseStatusCode)
	} else {
		cors.handleNormal(c)
	}

	if !cors.allowAllOrigins {
		c.Header("Access-Control-Allow-Origin", origin)
	}
}

func (cors *cors) validateWildcardOrigin(origin string) bool {
	for _, w := range cors.wildcardOrigins {
		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
			return true
		}
		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
			return true
		}
		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
			return true
		}
	}

	return false
}

func (cors *cors) validateOrigin(origin string) bool {
	if cors.allowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if len(cors.wildcardOrigins) > 0 && cors.validateWildcardOrigin(origin) {
		return true
	}
	if cors.allowOriginFunc != nil {
		return cors.allowOriginFunc(origin)
	}
	return false
}

func (cors *cors) handlePreflight(c *gin.Context) {
	header := c.Writer.Header()
	for key, value := range cors.preflightHeaders {
		header[key] = value
	}
}

func (cors *cors) handleNormal(c *gin.Context) {
	header := c.Writer.Header()
	for key, value := range cors.normalHeaders {
		header[key] = value
	}
}

// Default returns the location middleware with default configuration.
func Default() gin.HandlerFunc {
	config := DefaultConfig()
	config.AllowAllOrigins = true
	return New(config)
}

// New returns the location middleware with user-defined custom configuration.
func New(config Config) gin.HandlerFunc {
	cors := newCors(config)
	return func(c *gin.Context) {
		cors.applyCors(c)
	}
}
