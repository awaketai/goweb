package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/spf13/cast"
)

type IRequest interface {
	// 请求地址中的参数：foo.com?a=1&b=2&c[]=bar
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) any

	// 路由地址中带的参数 /book/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) any

	// form表单中带的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) any

	// json body
	BindJson(obj any) error

	// xml body
	BindXml(obj any) error

	// other
	GetRqwData() ([]byte, error)

	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header
	Headers() map[string][]string
	Header(key string) (string, bool)

	// cookie

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

var _ IRequest = new(Context)

const defaultMultipartMemory = 32 << 20 // 32mb

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}

	return map[string][]string(ctx.request.URL.Query())
}

func (ctx *Context) Query(key string) any {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}
	return nil
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryString(key, def string) (string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) Param(key string) any {
	if ctx.params == nil {
		return nil
	}
	if val, ok := ctx.params[key]; ok {
		return val
	}
	return nil
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}

	return map[string][]string(ctx.request.PostForm)
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormString(key, def string) (string, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals[0], true
	}
	return def, false
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	parmas := ctx.FormAll()
	if vals, ok := parmas[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}

	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (ctx *Context) Form(key string) any {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals
		}
	}
	return nil
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

func (ctx *Context) BindXml(obj any) error {
	if ctx.request == nil {
		return fmt.Errorf("ctx.request empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = xml.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) GetRqwData() ([]byte, error) {
	if ctx.request == nil {
		return nil, fmt.Errorf("ctx.request empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return nil, err
	}
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	return ctx.request.Host
}

func (ctx *Context) ClientIp() string {
	r := ctx.request
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

func (ctx *Context) Headers() map[string][]string {
	return ctx.request.Header
}

func (ctx *Context) Header(key string) (string, bool) {
	vals := ctx.request.Header.Values(key)
	if len(vals) == 0 {
		return "", false
	}
	return vals[0], true
}

func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookies := ctx.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}
