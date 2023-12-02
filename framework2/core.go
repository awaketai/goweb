package framework2

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router      map[string]*Tree
	middlewares []ControllerHandler
}

// 1.definition route map
// 2.register route
// 3.match route
// 4.implement ServerHTTP method

func NewCore() *Core {
	router := map[string]*Tree{}
	router[http.MethodGet] = NewTree()
	router[http.MethodPost] = NewTree()
	router[http.MethodPut] = NewTree()
	router[http.MethodDelete] = NewTree()
	core := &Core{router: router}
	return core
}

func (c *Core) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	//
	ctx := NewContext(w, r)
	handlers := c.FindRuteByRequest(r)
	if handlers == nil {
		ctx.JSON(404, "not found")
		return
	}
	ctx.SetHandlers(handlers)
	
	// call func
	if err := ctx.Next();err != nil {
		ctx.JSON(500,"not found")
		return
	}
}

func (c *Core) Use(middlewarese ...ControllerHandler){
	c.middlewares = append(c.middlewares, middlewarese...)
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) Get(url string, handler ...ControllerHandler) {
	allhandlers := append(c.middlewares,handler...)
	if err := c.router[http.MethodGet].AddRouter(url,allhandlers);err != nil {
		log.Fatal("add router err:",err)
	}
}

func (c *Core) Post(url string, handler ...ControllerHandler) {
	allhandlers := append(c.middlewares,handler...)
	if err := c.router[http.MethodPost].AddRouter(url,allhandlers);err != nil {
		log.Fatal("add router err:",err)
	}
}

func (c *Core) Put(url string, handler ...ControllerHandler) {
	allhandlers := append(c.middlewares,handler...)
	if err := c.router[http.MethodPut].AddRouter(url,allhandlers);err != nil {
		log.Fatal("add router err:",err)
	}
}

func (c *Core) Delete(url string, handler ...ControllerHandler) {
	allhandlers := append(c.middlewares,handler...)
	if err := c.router[http.MethodDelete].AddRouter(url,allhandlers);err != nil {
		log.Fatal("add router err:",err)
	}
}

func (c *Core) FindRuteByRequest(request *http.Request) []ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	// [method][uri]
	if methodHandlers, ok := c.router[strings.ToUpper(method)]; ok {
		return methodHandlers.FindHandler(uri)
	}

	return nil
}
