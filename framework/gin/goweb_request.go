package gin

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
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)
	DefaultQuery(key, def string) string

	// 路由地址中带的参数 /book/:id
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	DefaultParam(key string) any

	// form表单中带的参数
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string) any

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
	Header(key, val string)

	// cookie
	Cookies() map[string]string
	Cookie(key string) (string, error)
}

var _ IRequest = new(Context)

func (ctx *Context) QueryAll() map[string][]string {
	ctx.initQueryCache()
	if ctx.Request == nil {
		return map[string][]string{}
	}

	return map[string][]string(ctx.Request.URL.Query())
}

func (ctx *Context) DefaultQueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryString(key, def string) (string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) DefaultParam(key string) any {
	if val, ok := ctx.Params.Get(key); ok {
		return val
	}
	return nil
}

func (ctx *Context) DefaultParamInt(key string, def int) (int, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamString(key string, def string) (string, bool) {
	if val := ctx.DefaultParam(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

func (ctx *Context) FormAll() map[string][]string {
	ctx.initFormCache()
	return map[string][]string(ctx.formCache)
}

func (ctx *Context) DefaultFormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormString(key, def string) (string, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals[0], true
	}
	return def, false
}

func (ctx *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	parmas := ctx.FormAll()
	if vals, ok := parmas[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) DefaultFormFile(key string) (*multipart.FileHeader, error) {
	if ctx.Request.MultipartForm == nil {
		if err := ctx.Request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := ctx.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (ctx *Context) DefaultForm(key string) any {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals
		}
	}
	return nil
}

func (ctx *Context) BindJson(obj any) error {
	if ctx.Request == nil {
		return fmt.Errorf("ctx.request empty")
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) BindXml(obj any) error {
	if ctx.Request == nil {
		return fmt.Errorf("ctx.Request empty")
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = xml.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) GetRqwData() ([]byte, error) {
	if ctx.Request == nil {
		return nil, fmt.Errorf("ctx.Request empty")
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func (ctx *Context) Uri() string {
	return ctx.Request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) Host() string {
	return ctx.Request.Host
}

func (ctx *Context) ClientIp() string {
	r := ctx.Request
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
	return ctx.Request.Header
}

func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.Request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}
