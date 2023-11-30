package goweb

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	// handlers current request chain
	handlers []ControllerHandler
	// index arrive which node 
	index int
	isTimeout      bool
	writeMux       *sync.Mutex
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		// gurantee the first handler is the root handler
		index: -1,
		writeMux:       &sync.Mutex{},
	}
}

func (c *Context) WriteMux() *sync.Mutex {
	return c.writeMux
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) IsTimeout() bool {
	return c.isTimeout
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key any) any {
	return c.BaseContext().Value(key)
}

// Next next handler
func (c *Context) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		if err := c.handlers[c.index](c);err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) SetHandlers(handlers []ControllerHandler){
	c.handlers = handlers
}

func (c *Context) JSON(status int, obj any) error {
	if c.IsTimeout() {
		return nil
	}

	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)
	byts, err := json.Marshal(obj)
	if err != nil {
		c.responseWriter.WriteHeader(500)
		return err
	}
	_, err = c.responseWriter.Write(byts)
	if err != nil {
		c.responseWriter.WriteHeader(500)
		return err
	}
	return nil
}
