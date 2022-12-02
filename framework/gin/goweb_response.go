package gin

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/awaketai/goweb/framework/gin/internal/json"
)

type IResponse interface {
	IJson(obj any) IResponse
	IJsonp(obj any) IResponse
	IXml(obj any) IResponse
	IHtml(file string, obj any) IResponse
	IText(format string, values ...any) IResponse
	IRedirect(path string) IResponse
	ISetHeader(key, val string) IResponse
	ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	ISetStatus(code int) IResponse
	// 200 status
	ISetOkStatus() IResponse
}

var _ IResponse = new(Context)

func (ctx *Context) IJson(obj any) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/json")
	ctx.Writer.Write(byt)
	return ctx
}

func (ctx *Context) IJsonp(obj any) IResponse {
	// 获取请求参数callback
	callbackFunc, _ := ctx.DefaultQueryString("callback", "callback_function")
	ctx.ISetHeader("Content-Type", "application/javascript")
	// 转义
	callback := template.JSEscapeString(callbackFunc)
	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		ctx.ISetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		ctx.ISetStatus(http.StatusInternalServerError)
		return ctx
	}
	ret, err := json.Marshal(obj)
	if err != nil {
		ctx.ISetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.Writer.Write(ret)
	if err != nil {
		ctx.ISetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		ctx.ISetStatus(http.StatusInternalServerError)
		return ctx
	}

	return ctx
}

func (ctx *Context) IXml(obj any) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/html")
	ctx.Writer.Write(byt)
	return ctx

}

func (ctx *Context) IHtml(file string, obj any) IResponse {
	tpl, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	if err := tpl.Execute(ctx.Writer, obj); err != nil {
		return ctx
	}
	ctx.ISetHeader("Content-Type", "application/html")
	return ctx
}

func (ctx *Context) IText(format string, values ...any) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.ISetHeader("Content-Type", "application/text")
	ctx.Writer.Write([]byte(out))
	return ctx
}

func (ctx *Context) IRedirect(path string) IResponse {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *Context) ISetStatus(code int) IResponse {
	ctx.Writer.WriteHeader(code)
	return ctx
}

func (ctx *Context) ISetHeader(key, val string) IResponse {
	ctx.Writer.Header().Add(key, val)
	return ctx
}

func (ctx *Context) ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}

func (ctx *Context) ISetOkStatus() IResponse {
	ctx.Writer.WriteHeader(http.StatusOK)
	return ctx
}
