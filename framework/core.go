package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router      map[string]*Trie
	middlewares []ControllerHandler
}

func NewCore() *Core {
	// write
	router := make(map[string]*Trie)
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{router: router}
}

func (c *Core) RegisterRouter(url string, handler ControllerHandler) {
	// c.router[url] = handler
}

func (c *Core) FindRouteByRequest(req *http.Request) []ControllerHandler {
	uri := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)
	log.Printf("Method:%s uri:%s", req.Method, req.URL.Path)
	// first level
	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.FindHandler(uri)
	}

	return nil
}

func (c *Core) FindRouteNodeByRequest(req *http.Request) *node {
	uri := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)
	log.Printf("Method:%s uri:%s", req.Method, req.URL.Path)
	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	ctx := NewContext(r, w)
	// find routes
	node := c.FindRouteNodeByRequest(r)
	if node == nil {

		ctx.SetStatus(http.StatusNotFound).Json("not found")
		return
	}

	ctx.SetHandlers(node.handlers)
	params := node.parseParamsFromEndNode(r.URL.Path)
	ctx.SetParams(params)
	// invoke
	if err := ctx.Next(); err != nil {
		ctx.SetStatus(http.StatusInternalServerError).Json("inner error")
		return
	}
}

// registe middleware
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "GET", url)
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "POST", url)

}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "PUT", url)

}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "DELETE", url)

}
