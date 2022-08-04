package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Trie
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

func (c *Core) FindRouteByRequest(req *http.Request) ControllerHandler {
	uri := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)
	log.Printf("Method:%s uri:%s", req.Method, req.URL.Path)
	// first level
	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.FindHandler(uri)
	}

	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	ctx := NewContext(r, w)
	// find route
	handler := c.FindRouteByRequest(r)
	if handler == nil {
		ctx.Json(http.StatusNotFound, "not found")
		return
	}

	// invoke
	if err := handler(ctx); err != nil {
		ctx.Json(http.StatusInternalServerError, "inner error")
		return
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	if err := c.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "GET", url)
}

func (c *Core) Post(url string, handler ControllerHandler) {
	if err := c.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "POST", url)

}

func (c *Core) Put(url string, handler ControllerHandler) {
	if err := c.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "PUT", url)

}

func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add route error:", err)
	}
	log.Printf("add route Method:%s uri:%s", "DELETE", url)

}
