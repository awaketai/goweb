package kernel

import (
	"net/http"

	"github.com/awaketai/goweb/framework/gin"
)

type WebKernelService struct {
	engine *gin.Engine
}

func NewWebKernelService(params ...any) (any, error) {
	if len(params) == 0 {
		panic("parmas error")
	}
	httpEngine := params[0].(*gin.Engine)
	return &WebKernelService{engine: httpEngine}, nil
}

func (service *WebKernelService) HttpEngine() http.Handler {
	return service.engine
}
