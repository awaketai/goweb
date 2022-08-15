package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	Json(obj any) IResponse
	Jsonp(obj any) IResponse
	Xml(obj any) IResponse
	Html(file string, obj any) IResponse
	Text(format string, values ...any) IResponse
	Redirect(path string) IResponse
	SetHeader(key, val string) IResponse
	SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	SetStatus(code int) IResponse
	// 200 status
	SetOkStatus() IResponse
}

var _ IResponse = new(Context)

func (ctx *Context) Json(obj any) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.responseWriter.Write(byt)
	return ctx
}

func (ctx *Context) Jsonp(obj any) IResponse {
	// 获取请求参数callback
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Content-Type", "application/javascript")
	// 转义
	callback := template.JSEscapeString(callbackFunc)
	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return ctx
	}
	ret, err := json.Marshal(obj)
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.responseWriter.Write(ret)
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return ctx
	}
	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return ctx
	}

	return ctx
}

func (ctx *Context) Xml(obj any) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/html")
	ctx.responseWriter.Write(byt)
	return ctx

}

func (ctx *Context) Html(file string, obj any) IResponse {
	tpl, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	if err := tpl.Execute(ctx.responseWriter, obj); err != nil {
		return ctx
	}
	ctx.SetHeader("Content-Type", "application/html")
	return ctx
}

func (ctx *Context) Text(format string, values ...any) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/text")
	ctx.responseWriter.Write([]byte(out))
	return ctx
}

func (ctx *Context) Redirect(path string) IResponse {
	http.Redirect(ctx.responseWriter, ctx.request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *Context) SetStatus(code int) IResponse {
	ctx.responseWriter.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetHeader(key, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}

func (ctx *Context) SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.responseWriter, &http.Cookie{
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

func (ctx *Context) SetOkStatus() IResponse {
	ctx.responseWriter.WriteHeader(http.StatusOK)
	return ctx
}
