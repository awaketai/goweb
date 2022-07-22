package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	hasTimeout     bool // 是否超时标记位
	writerMux      *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
	}
}

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.request.Context().Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.request.Context().Value(key)
}

func (ctx *Context) Deadline() (time.Time, bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}

	return map[string][]string(ctx.request.URL.Query())
}

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len == 0 {
			return def
		}

		intval, err := strconv.Atoi(vals[len-1])
		if err != nil {
			return def
		}
		return intval
	}
	return def
}

func (ctx *Context) QueryString(key, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	parmas := ctx.QueryAll()
	if vals, ok := parmas[key]; ok {
		return vals
	}

	return def
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}

	return map[string][]string(ctx.request.PostForm)
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len == 0 {
			return def
		}

		intval, err := strconv.Atoi(vals[len-1])
		if err != nil {
			return def
		}
		return intval
	}
	return def
}

func (ctx *Context) FormString(key, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	parmas := ctx.FormAll()
	if vals, ok := parmas[key]; ok {
		return vals
	}

	return def
}

func (ctx *Context) BindJson(obj any) error {
	if ctx.request == nil {
		return fmt.Errorf("ctx.request empty")
	}

	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}

	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) Json(status int, obj any) error {
	if ctx.HasTimeout() {
		return nil
	}

	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}

	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj any, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj any) error {
	return nil
}
