package goweb

import (
	"net/http"
	"strings"
)

type Core struct {
	router map[string]map[string]ControllerHandler
}

// 1.definition route map
// 2.register route
// 3.match route
// 4.implement ServerHTTP method

func NewCore() *Core {
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}
	router := map[string]map[string]ControllerHandler{}
	router[http.MethodGet] = getRouter
	router[http.MethodPost] = postRouter
	router[http.MethodPut] = putRouter
	router[http.MethodDelete] = deleteRouter
	return &Core{router: router}
}

func (c *Core) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	//
	ctx := NewContext(w, r)
	router := c.FindRuteByRequest(r)
	if router == nil {
		ctx.JSON(404, "not found")
		return
	}
	// call func
	if err := router(nil); err != nil {
		ctx.JSON(500, "inner error")
		return
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[http.MethodGet][strings.ToUpper(url)] = handler
}

func (c *Core) Post(url string, handler ControllerHandler) {
	c.router[http.MethodPost][strings.ToUpper(url)] = handler
}

func (c *Core) Put(url string, handler ControllerHandler) {
	c.router[http.MethodPut][strings.ToUpper(url)] = handler
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	c.router[http.MethodDelete][strings.ToUpper(url)] = handler
}

func (c *Core) FindRuteByRequest(request *http.Request) ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	// [method][uri]
	if methodHandlers, ok := c.router[strings.ToUpper(method)]; ok {
		if handler, ok := methodHandlers[strings.ToUpper(uri)]; ok {
			return handler
		}
	}

	return nil
}
